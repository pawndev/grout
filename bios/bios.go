package bios

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"grout/cfw"
	"grout/constants"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func SaveFile(biosFile constants.BIOSFile, platformFSSlug string, data []byte) error {
	filePaths := cfw.GetBIOSFilePaths(biosFile.RelativePath, platformFSSlug)

	for _, filePath := range filePaths {
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	return nil
}

func VerifyFileMD5(data []byte, expectedMD5 string) (bool, string) {
	if expectedMD5 == "" {
		// No MD5 hash to verify against
		return true, ""
	}

	hash := md5.Sum(data)
	actualMD5 := hex.EncodeToString(hash[:])

	return actualMD5 == expectedMD5, actualMD5
}

func GetFileInfo(biosFile constants.BIOSFile, platformFSSlug string) (exists bool, size int64, md5Hash string, err error) {
	filePaths := cfw.GetBIOSFilePaths(biosFile.RelativePath, platformFSSlug)

	for _, filePath := range filePaths {
		info, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return false, 0, "", err
		}

		file, err := os.Open(filePath)
		if err != nil {
			return true, info.Size(), "", err
		}
		defer file.Close()

		hash := md5.New()
		if _, err := io.Copy(hash, file); err != nil {
			return true, info.Size(), "", err
		}

		md5Hash = hex.EncodeToString(hash.Sum(nil))

		return true, info.Size(), md5Hash, nil
	}

	return false, 0, "", nil
}

func GetFilesForPlatform(platformFSSlug string) []constants.BIOSFile {
	var biosFiles []constants.BIOSFile

	coreNames, ok := constants.PlatformToLibretroCores[platformFSSlug]
	if !ok {
		return biosFiles
	}

	seen := make(map[string]bool)
	for _, coreName := range coreNames {
		normalizedCoreName := strings.TrimSuffix(coreName, "_libretro")
		coreInfo, ok := constants.LibretroCoreToBIOS[normalizedCoreName]
		if !ok {
			continue
		}

		for _, file := range coreInfo.Files {
			if !seen[file.FileName] {
				biosFiles = append(biosFiles, file)
				seen[file.FileName] = true
			}
		}
	}

	return biosFiles
}

type Status string

const (
	StatusMissing        Status = "missing"
	StatusValid          Status = "valid"
	StatusInvalidHash    Status = "invalid_hash"
	StatusNoHashToVerify Status = "no_hash"
)

type FileStatus struct {
	File        constants.BIOSFile
	Status      Status
	Exists      bool
	Size        int64
	ActualMD5   string
	ExpectedMD5 string
}

func CheckFileStatus(biosFile constants.BIOSFile, platformFSSlug string) FileStatus {
	status := FileStatus{
		File:        biosFile,
		ExpectedMD5: biosFile.MD5Hash,
	}

	exists, size, actualMD5, err := GetFileInfo(biosFile, platformFSSlug)
	if err != nil {
		status.Status = StatusMissing
		return status
	}

	status.Exists = exists
	status.Size = size
	status.ActualMD5 = actualMD5

	if !exists {
		status.Status = StatusMissing
		return status
	}

	if biosFile.MD5Hash == "" {
		status.Status = StatusNoHashToVerify
		return status
	}

	if actualMD5 == biosFile.MD5Hash {
		status.Status = StatusValid
	} else {
		status.Status = StatusInvalidHash
	}

	return status
}
