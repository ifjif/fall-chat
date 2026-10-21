package errs

import "github.com/gofiber/fiber/v3"

var (
	FileErrChunkCountInvalid = NewBizError(4400, fiber.StatusBadRequest, "文件chunk数不正确")
	FileErrTooLarge          = NewBizError(4401, fiber.StatusBadRequest, "文件太大请分片")
)
