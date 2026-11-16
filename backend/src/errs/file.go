package errs

import "github.com/gofiber/fiber/v3"

var (
	FileErrChunkCountInvalid = NewBizError(FileStartCode, fiber.StatusBadRequest, "文件chunk数不正确")
	FileErrTooLarge          = NewBizError(FileStartCode+1, fiber.StatusBadRequest, "文件太大请分片")
)
