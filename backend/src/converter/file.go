package converter

import (
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
)

func ToFileConfigResponse(m *model.FileConfig) *response.FileConfig {
	return &response.FileConfig{
		ID:        m.PublicID,
		MaxSize:   m.MaxSize,
		ChunkSize: m.ChunkSize,
	}
}

func ToFileUploadChunkCommand(req *request.FileUploadChunk) *command.FileUploadChunk {
	return &command.FileUploadChunk{
		FileID: req.FileID,
		Index:  req.Index,
		Chunk:  req.Chunk,
	}
}

func ToFileUploadMergeCommand(req *request.FileUploadMerge) *command.FileUploadMerge {
	return &command.FileUploadMerge{
		FileID:     req.FileID,
		Filename:   req.Filename,
		TotalChunk: req.TotalChunk,
		Chunks:     req.Chunks,
	}
}

func ToFileInstanceCheckQuery(req *request.FileInstanceCheck) *query.FileInstanceCheck {
	return &query.FileInstanceCheck{
		FileID: req.FileID,
	}
}

func ToFileChunkStatusQuery(req *request.FileChunkStatus) *query.FileChunkStatus {
	return &query.FileChunkStatus{
		FileID: req.FileID,
	}
}

func ToFileUploadCommand(req *request.FileUpload) *command.FileUpload {
	return &command.FileUpload{
		FileID: req.FileID,
		File:   req.File,
	}
}

func ToFileFetchQuery(req *request.FileFetch) *query.FileFetch {
	return &query.FileFetch{
		File: req.File,
	}
}
