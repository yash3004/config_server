package configurations

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConfigManager handles configuration file operations
type ConfigManager struct {
	db        *mongo.Database
	useFile   bool
	configDir string
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(db *mongo.Database, useFile bool, configDir string) *ConfigManager {
	if configDir == "" {
		configDir = "configurations"
	}
	return &ConfigManager{
		db:        db,
		useFile:   useFile,
		configDir: configDir,
	}
}

// AddConfig adds a new configuration file to GridFS or local filesystem
func (cm *ConfigManager) AddConfig(ctx context.Context, userID, filename string, fileType int, data []byte) error {
	if cm.useFile {
		return cm.addConfigToFile(userID, filename, data)
	}

	bucket, err := gridfs.NewBucket(cm.db)
	if err != nil {
		return err
	}

	// Create metadata for the file
	metadata := bson.M{
		"userID":   userID,
		"fileType": fileType,
		"created":  time.Now(),
		"updated":  time.Now(),
	}

	// Check if file already exists
	if cm.fileExists(ctx, userID, filename) {
		return errors.New("file already exists")
	}

	uploadOpts := options.GridFSUpload().SetMetadata(metadata)
	uploadStream, err := bucket.OpenUploadStream(filename, uploadOpts)
	if err != nil {
		return err
	}
	defer uploadStream.Close()

	_, err = uploadStream.Write(data)
	return err
}

// addConfigToFile adds a new configuration file to the local filesystem
func (cm *ConfigManager) addConfigToFile(userID, filename string, data []byte) error {
	// Create path to user's configuration directory
	userDir := filepath.Join(cm.configDir, userID)

	// Create user directory if it doesn't exist
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(userDir, filename)

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		return errors.New("file already exists")
	}

	// Write file
	return ioutil.WriteFile(filePath, data, 0644)
}

// UpdateConfig updates an existing configuration file in GridFS or local filesystem
func (cm *ConfigManager) UpdateConfig(ctx context.Context, userID, filename string, fileType int, data []byte) error {
	if err := cm.DeleteConfig(ctx, userID, filename); err != nil {
		return err
	}

	return cm.AddConfig(ctx, userID, filename, fileType, data)
}

func (cm *ConfigManager) DeleteConfig(ctx context.Context, userID, filename string) error {
	if cm.useFile {
		return cm.deleteConfigFromFile(userID, filename)
	}

	bucket, err := gridfs.NewBucket(cm.db)
	if err != nil {
		return err
	}

	fileID, err := cm.getFileID(ctx, userID, filename)
	if err != nil {
		return err
	}

	return bucket.Delete(fileID)
}

func (cm *ConfigManager) deleteConfigFromFile(userID, filename string) error {
	userDir := filepath.Join(cm.configDir, userID)
	filePath := filepath.Join(userDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errors.New("file not found")
	}

	return os.Remove(filePath)
}

func (cm *ConfigManager) GetConfig(ctx context.Context, userID, filename string) ([]byte, int, error) {
	if cm.useFile {
		return cm.getConfigFromFile(userID, filename)
	}

	bucket, err := gridfs.NewBucket(cm.db)
	if err != nil {
		return nil, 0, err
	}

	fileID, err := cm.getFileID(ctx, userID, filename)
	if err != nil {
		return nil, 0, err
	}

	var file struct {
		Metadata struct {
			FileType int    `bson:"fileType"`
			UserID   string `bson:"userID"`
		} `bson:"metadata"`
	}

	err = cm.db.Collection("fs.files").FindOne(ctx, bson.M{"_id": fileID}).Decode(&file)
	if err != nil {
		return nil, 0, err
	}

	if file.Metadata.UserID != userID {
		return nil, 0, errors.New("unauthorized access to file")
	}

	downloadStream, err := bucket.OpenDownloadStream(fileID)
	if err != nil {
		return nil, 0, err
	}
	defer downloadStream.Close()

	data, err := io.ReadAll(downloadStream)
	if err != nil {
		return nil, 0, err
	}

	return data, file.Metadata.FileType, nil
}

func (cm *ConfigManager) getConfigFromFile(userID, filename string) ([]byte, int, error) {
	userDir := filepath.Join(cm.configDir, userID)
	filePath := filepath.Join(userDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, 0, errors.New("file not found")
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, 0, err
	}

	fileType := determineFileType(filepath.Ext(filename))

	return data, fileType, nil
}

func (cm *ConfigManager) fileExists(ctx context.Context, userID, filename string) bool {
	_, err := cm.getFileID(ctx, userID, filename)
	return err == nil
}

func (cm *ConfigManager) getFileID(ctx context.Context, userID, filename string) (primitive.ObjectID, error) {
	var file struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	filter := bson.M{
		"filename":        filename,
		"metadata.userID": userID,
	}

	err := cm.db.Collection("fs.files").FindOne(ctx, filter).Decode(&file)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return file.ID, nil
}

func determineFileType(ext string) int {
	switch ext {
	case ".txt":
		return 1
	case ".csv":
		return 2
	case ".json":
		return 3
	case ".xml":
		return 4
	case ".yaml", ".yml":
		return 5
	default:
		return 1
	}
}
