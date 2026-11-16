package controller

import (
	"fmt"
	"net/url"
	"path"
	"strconv"
	"zzc/fall-chat/src/bind"
	"zzc/fall-chat/src/converter"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/filestore"
	"zzc/fall-chat/src/keys/ctxkey"
	"zzc/fall-chat/src/middleware"
	"zzc/fall-chat/src/resp"
	"zzc/fall-chat/src/service"
	"zzc/fall-chat/src/utils"

	"github.com/gofiber/fiber/v3"
)

type File struct {
	Base
	fileService *service.File
}

func NewFile(name string) *File {
	bs := Base{name: name}
	return &File{
		Base: bs,
	}
}

func (fc *File) Inject() {
	fs := service.GlobalService[service.FILE_SERVICE_NAME].(*service.File)
	fc.fileService = fs
}

func (fc *File) SetRoute(r fiber.Router) {
	group := r.Group("/file")

	authGroup := group.Group("")
	authGroup.Use(middleware.AuthMiddleware)

	authGroup.Get("/config", fc.config)

	authGroup.Get("/fetch", fc.fetch)
	authGroup.Get("/accessUrl", fc.accessURL)
	authGroup.Post("/upload", fc.upload)
	authGroup.Post("/chunk", fc.chunk)
	authGroup.Post("/chunkStatus", fc.chunkStatus)
	authGroup.Post("/merge", fc.merge)
	authGroup.Post("/instanceCheck", fc.instanceCheck)
}

// /file/config 	需auth
func (fc *File) config(c fiber.Ctx) error {
	result, err := fc.fileService.Config(c.Context())
	if err != nil {
		return err
	}

	fcRespone := converter.ToFileConfigResponse(result)

	return resp.
		New(c).
		OK().
		JSON(fcRespone)
}

// /file/fetch 需auth
func (fc *File) fetch(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Query[request.FileFetch](c)
	if err != nil {
		return fc.Wrap("body bind", err,
			"username", username,
		)
	}

	qu := converter.ToFileFetchQuery(req)
	qu.ID = id
	qu.Username = username

	result, err := fc.fileService.Fetch(qu)
	if err != nil {
		return err
	}

	size := result.Stat.Size
	filename := path.Base(result.Stat.Key)

	disposition := fmt.Sprintf(`attachment;filename="%s";filename*=UTF-8''%s`, filename, url.PathEscape(filename))
	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Disposition", disposition)
	c.Set("Content-Length", strconv.FormatInt(size, 10))
	c.Set("Cache-Control", "private,max-age=300")

	return c.SendStream(result.Object, int(size))
}

// /file/upload 需auth
// 小于5MB 走它
func (fc *File) upload(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)
	req, err := bind.Body[request.FileUpload](c)
	if err != nil {
		return fc.Wrap("body bind", err, "username", username)
	}

	filesize := req.File.Size
	if filesize > filestore.ChunkSize {
		return fc.Wrap("upload", errs.FileErrTooLarge, "username", username)
	}

	cmd := converter.ToFileUploadCommand(req)
	cmd.UserID = id
	cmd.Username = username

	path, err := fc.fileService.Upload(cmd)
	if err != nil {
		return err
	}
	return resp.
		New(c).
		OK().
		Message("上传成功").
		JSON(path)
}

// /file/chunk 需auth
func (fc *File) chunk(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)
	req, err := bind.Form[request.FileUploadChunk](c)
	if err != nil {
		return fc.Wrap("form绑定", err)
	}

	cmd := converter.ToFileUploadChunkCommand(req)
	cmd.UserID = id
	cmd.Username = username

	err = fc.fileService.UploadChunk(cmd)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(nil)
}

func (fc *File) chunkStatus(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	usernmae := ctxkey.GetUsername(c)

	req, err := bind.Body[request.FileChunkStatus](c)
	if err != nil {
		return fc.Wrap("body bind", err, "username", usernmae)
	}

	qu := converter.ToFileChunkStatusQuery(req)
	qu.UserID = id
	qu.Username = usernmae

	result, err := fc.fileService.ChunkStatus(qu)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(result)
}

// /file/merge 需auth
func (fc *File) merge(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.FileUploadMerge](c)
	if err != nil {
		return fc.Wrap("body绑定", err,
			"username", username,
		)
	}

	cmd := converter.ToFileUploadMergeCommand(req)
	cmd.UserID = id
	cmd.Username = username

	filename, err := fc.fileService.UploadMerge(cmd)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		Message("上传成功").
		JSON(filename)
}

// /file/instanceCheck 需auth
func (fc *File) instanceCheck(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)
	req, err := bind.Body[request.FileInstanceCheck](c)
	if err != nil {
		return fc.Wrap("Body bind", err,
			"username", username,
		)
	}

	qu := converter.ToFileInstanceCheckQuery(req)
	qu.UserID = id
	qu.Username = username

	fileURL, err := fc.fileService.InstanceCheck(qu)
	if err != nil {
		return err
	}

	if fileURL == "" {
		return resp.
			New(c).
			Err(fiber.StatusNotFound).
			Failure()
	}

	return resp.
		New(c).
		OK().JSON(fileURL)
}

// /file/accessUrl 需auth
func (fc *File) accessURL(c fiber.Ctx) error {
	id := ctxkey.GetUserPubID(c)
	username := ctxkey.GetUsername(c)
	filename := c.Query("file")
	if filename == "" {
		return fc.Wrap("accessUrl", bind.ErrBind, "username", username)
	}
	// 生成临时token
	token, _, err := utils.GenerateFileAccessToken(id, username, filename)
	if err != nil {
		return fc.Wrap("accessUrl", errs.CommonErrJWTTokenGeneration.Meger(err), "username", username)
	}

	url := fmt.Sprintf("/api/v1/file/fetch?file=%s&token=%s", url.QueryEscape(filename), url.QueryEscape("Bearer "+token))

	c.Set("Cache-Control", "private, max-age=300")

	return resp.
		New(c).
		OK().
		JSON(url)
}
