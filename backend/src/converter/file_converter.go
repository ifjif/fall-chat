package converter

import (
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
)

func ToFileConfigResponse(m *model.FileConfig) *response.FileConfigResponse {
	return &response.FileConfigResponse{
		ID:        m.PublicID,
		MaxSize:   m.MaxSize,
		ChunkSize: m.ChunkSize,
	}
}

func ToFileUploadChunkCommand(req *request.FileUploadChunkRequest) *command.FileUploadChunkCommand {
	return &command.FileUploadChunkCommand{
		FileID: req.FileID,
		Index:  req.Index,
		Chunk:  req.Chunk,
	}
}

func ToFileUploadMergeCommand(req *request.FileUploadMergeRequest) *command.FileUploadMergeCommand {
	return &command.FileUploadMergeCommand{
		FileID:     req.FileID,
		Filename:   req.Filename,
		TotalChunk: req.TotalChunk,
		Chunks:     req.Chunks,
	}
}

func ToFileInstanceCheckQuery(req *request.FileInstanceCheckRequest) *query.FileInstanceCheckQuery {
	return &query.FileInstanceCheckQuery{
		FileID: req.FileID,
	}
}

func ToFileChunkStatusQuery(req *request.FileChunkStatusRequest) *query.FileChunkStatusQuery {
	return &query.FileChunkStatusQuery{
		FileID: req.FileID,
	}
}

func ToFileUploadCommand(req *request.FileUploadRequest) *command.FileUploadCommand {
	return &command.FileUploadCommand{
		FileID: req.FileID,
		File:   req.File,
	}
}

func ToFileFetchQuery(req *request.FileFetchRequest) *query.FileFetchQuery {
	return &query.FileFetchQuery{
		File: req.File,
	}
}
