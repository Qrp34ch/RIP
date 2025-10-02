package repository

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"strings"
	"time"

	"lab1/internal/app/ds"
)

func (r *Repository) GetReactions() ([]ds.Reaction, error) {
	var reactions []ds.Reaction
	err := r.db.Where("is_delete = ?", false).Find(&reactions).Error
	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
	if err != nil {
		return nil, err
	}
	if len(reactions) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return reactions, nil
}

func (r *Repository) GetReaction(id int) (ds.Reaction, error) {
	reaction := ds.Reaction{}
	err := r.db.Where("id = ? AND is_delete = ?", id, false).First(&reaction).Error
	if err != nil {
		return ds.Reaction{}, err
	}
	return reaction, nil
}

func (r *Repository) GetReactionsByTitle(title string) ([]ds.Reaction, error) {
	var reactions []ds.Reaction
	err := r.db.Where("title ILIKE ? AND is_delete = ?", "%"+title+"%", false).Find(&reactions).Error
	if err != nil {
		return nil, err
	}
	return reactions, nil
}

func (r *Repository) GetReactionsInSynthesis() int64 {
	var synthesisID uint
	var count int64
	creatorID := 1
	// пока что мы захардкодили id создателя заявки, в последующем вы сделаете авторизацию и будете получать его из JWT

	err := r.db.Model(&ds.Synthesis{}).Where("creator_id = ? AND status = ?", creatorID, "черновик").Select("id").First(&synthesisID).Error
	if err != nil {
		return 0
	}

	err = r.db.Model(&ds.SynthesisReaction{}).Where("synthesis_id = ?", synthesisID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in lists_reactions:", err)
	}

	return count
}

func (r *Repository) GetSynthesis(synthesisID uint) ([]ds.Reaction, error) {
	var reactionID uint
	var reaction ds.Reaction
	var synthesisReaction []ds.SynthesisReaction
	var result []ds.Reaction
	err := r.db.Where("synthesis_id = ?", synthesisID).Find(&synthesisReaction).Error
	if err != nil {
		return nil, err
	}

	for _, mm := range synthesisReaction {
		reactionID = mm.ReactionID
		reaction, err = r.GetReaction(int(reactionID))
		//if err != nil {
		//	return nil, err
		//}
		if err != nil {
			logrus.Warnf("Reaction %d not found or deleted, skipping", reactionID)
			continue
		}
		result = append(result, reaction)
	}

	return result, nil
}

func (r *Repository) FindUserSynthesis(userID uint) uint {
	var synthesisID uint
	err := r.db.Model(&ds.Synthesis{}).Where("creator_id = ? AND status = ?", userID, "черновик").Select("id").First(&synthesisID).Error
	if err != nil {
		return 0
	}
	return synthesisID
}

func (r *Repository) AddReactionInSynthesis(id uint) error {
	userID := 1
	moderatorID := 2
	var synthesisID uint
	var count int64

	err := r.db.Model(&ds.Synthesis{}).Where("creator_id = ? AND status = ?", userID, "черновик").Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		newSynthesis := ds.Synthesis{
			Status:      "черновик",
			DateCreate:  time.Now(),
			DateUpdate:  time.Now(),
			CreatorID:   uint(userID),
			ModeratorID: uint(moderatorID),
		}
		err := r.db.Create(&newSynthesis).Error
		if err != nil {
			return err
		}
	}

	err = r.db.Model(&ds.Synthesis{}).Where("creator_id = ? AND status = ?", userID, "черновик").Select("id").First(&synthesisID).Error
	if err != nil {
		return err
	}

	newSynthesisReaction := ds.SynthesisReaction{
		SynthesisID: synthesisID,
		ReactionID:  id,
	}

	err = r.db.Create(&newSynthesisReaction).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *Repository) RemoveSynthesis(id uint) error {
	deleteQuery := "UPDATE syntheses SET status = $1, date_finish = $2, date_update = $3 WHERE id = $4"
	r.db.Exec(deleteQuery, "удалён", time.Now(), time.Now(), id)
	return nil
}

func (r *Repository) GetUserNameByID(userID uint) string {
	var userName string
	r.db.Model(&ds.Users{}).Where("id = ?", userID).Select("fio").First(&userName)
	return userName
}

func (r *Repository) GetDateUpdate(synthesisID uint) string {
	var dateUpdateTime time.Time
	var dateUpdate string
	//не каждый человек знает песни Газманова так что
	//а я я ясные дни забираю себе а я хмурые дни возвращаю судьбе <3 ВЛАДИК ЧИКАНЧИ ПРИШЕЛ ЗА ШЛЯПОЙ В МАГАЗИН (ОНА ЕМУ НЕ ПОДОШЛА)
	r.db.Model(&ds.Synthesis{}).Where("id = ?", synthesisID).Select("date_update").First(&dateUpdateTime)
	dateUpdate = dateUpdateTime.Format("02.01.2006 15:04:05")
	return dateUpdate
}

func (r *Repository) SynthesisStatusById(synthesisID uint) (string, error) {
	var SynthesisStatus string
	err := r.db.Model(&ds.Synthesis{}).Where("id = ?", synthesisID).Select("status").First(&SynthesisStatus).Error
	if err != nil {
		return "", err
	}
	return SynthesisStatus, err
}

func (r *Repository) AddReaction(reaction *ds.Reaction) error {
	if reaction.Title == "" {
		return fmt.Errorf("название реакции обязательно")
	}

	//err := r.db.Select(
	//	"Title", "Details", "IsDelete", "StartingMaterial", "DensitySM",
	//	"VolumeSM", "MolarMassSM", "ResultMaterial", "DensityRM", "VolumeRM", "MolarMassRM",
	//).Create(reaction).Error
	err := r.db.Model(&ds.Reaction{}).Create(map[string]interface{}{
		"title":             reaction.Title,
		"details":           reaction.Details,
		"is_delete":         reaction.IsDelete,
		"starting_material": reaction.StartingMaterial,
		"density_sm":        reaction.DensitySM,
		"molar_mass_sm":     reaction.MolarMassSM,
		"result_material":   reaction.ResultMaterial,
		"density_rm":        reaction.DensityRM,
		"molar_mass_rm":     reaction.MolarMassRM,
	}).Error
	if err != nil {
		return fmt.Errorf("ошибка при создании реакции: %w", err)
	}
	return nil
}

func (r *Repository) ChangeReaction(id uint, reactionData *ds.Reaction) error {

	var reaction ds.Reaction
	err := r.db.Where("id = ? AND is_delete = false", id).First(&reaction).Error
	if err != nil {
		return fmt.Errorf("реакция с ID %d не найдена", id)
	}

	updReaction := map[string]interface{}{
		"title":             reactionData.Title,
		"src":               reactionData.Src,
		"src_ur":            reactionData.SrcUr,
		"details":           reactionData.Details,
		"is_delete":         reactionData.IsDelete,
		"starting_material": reactionData.StartingMaterial,
		"density_sm":        reactionData.DensitySM,
		"molar_mass_sm":     reactionData.MolarMassSM,
		"result_material":   reactionData.ResultMaterial,
		"density_rm":        reactionData.DensityRM,
		"molar_mass_rm":     reactionData.MolarMassRM,
	}

	for key, value := range updReaction {
		if value == "" || value == nil {
			delete(updReaction, key)
		}
	}

	err = r.db.Model(&ds.Reaction{}).Where("id = ?", id).Updates(updReaction).Error
	if err != nil {
		return fmt.Errorf("ошибка при обновлении реакции: %w", err)
	}

	return nil
}

func (r *Repository) DeleteReaction(id uint) error {
	var reaction ds.Reaction
	err := r.db.Where("id = ?", id).First(&reaction).Error
	if err != nil {
		return fmt.Errorf("реакция с ID %d не найдена: %w", id, err)
	}

	if reaction.Src != "" {
		if err := r.DeleteReactionImage(reaction.Src); err != nil {
			logrus.Errorf("Не удалось удалить изображение для реакции %d: %v", id, err)
		}
	}
	if reaction.SrcUr != "" {
		if err := r.DeleteReactionImage(reaction.SrcUr); err != nil {
			logrus.Errorf("Не удалось удалить изображение для реакции %d: %v", id, err)
		}
	}

	err = r.db.Model(&ds.Reaction{}).Where("id = ?", id).UpdateColumn("is_delete", true).Error
	fmt.Println(id)
	if err != nil {
		return fmt.Errorf("Ошибка при удалении реакции с id %d: %w", id, err)
	}

	if err := r.CleanupDeletedReactionsFromSyntheses(); err != nil {
		logrus.Warnf("Failed to cleanup deleted reactions: %v", err)
	}

	return nil
}

func (r *Repository) DeleteReactionImage(src string) error {
	if src == "" {
		logrus.Info("Empty image source, skipping deletion")
		return nil
	}
	objectName := r.extractObjectName(src)
	logrus.Infof("Extracted object name: '%s' from source: '%s'", objectName, src)
	if objectName == "" {
		logrus.Warnf("Could not extract object name from src: %s", src)
		return nil
	}
	logrus.Infof("Attempting to delete object: '%s' from bucket: '%s'", objectName, r.bucketName)

	_, err := r.minioClient.StatObject(context.Background(), r.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		logrus.Warnf("Object '%s' not found in MinIO: %v", objectName, err)
		return nil
	}

	err = r.minioClient.RemoveObject(context.Background(), r.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		logrus.Errorf("Failed to delete object '%s' from MinIO: %v", objectName, err)
		return fmt.Errorf("ошибка при удалении изображения из MinIO: %w", err)
	}

	logrus.Infof("Successfully deleted object: '%s' from bucket: '%s'", objectName, r.bucketName)
	return nil
}

func (r *Repository) extractObjectName(src string) string {
	logrus.Infof("Processing source URL: %s", src)
	if strings.Contains(src, "?") {
		src = strings.Split(src, "?")[0]
	}
	if strings.Contains(src, "http://localhost:9000/") {
		parts := strings.Split(src, "/")
		for i, part := range parts {
			if part == r.bucketName && i+1 < len(parts) {
				objectPath := strings.Join(parts[i+1:], "/")
				logrus.Infof("Extracted object path: %s", objectPath)
				return objectPath
			}
		}
	}

	if !strings.Contains(src, "/") {
		return "img/" + src
	}

	logrus.Infof("Using as-is: %s", src)
	return src
}

func (r *Repository) CleanupDeletedReactionsFromSyntheses() error {
	// Находим все удаленные реакции
	var deletedReactions []ds.Reaction
	err := r.db.Where("is_delete = ?", true).Find(&deletedReactions).Error
	if err != nil {
		return err
	}

	// Удаляем их из всех синтезов
	for _, reaction := range deletedReactions {
		err = r.db.Where("reaction_id = ?", reaction.ID).Delete(&ds.SynthesisReaction{}).Error
		if err != nil {
			logrus.Warnf("Failed to remove reaction %d from syntheses: %v", reaction.ID, err)
		} else {
			logrus.Infof("Removed deleted reaction %d from all syntheses", reaction.ID)
		}
	}

	return nil
}

func (r *Repository) UploadReactionImage(id uint, fileHeader *multipart.FileHeader) error {
	var reaction ds.Reaction
	err := r.db.Where("id = ? AND is_delete = false", id).First(&reaction).Error
	if err != nil {
		return fmt.Errorf("реакция с ID %d не найдена", id)
	}

	// Удаляем старое изображение если есть
	if reaction.Src != "" {
		if err := r.DeleteReactionImage(reaction.Src); err != nil {
			logrus.Errorf("Не удалось удалить старое изображение: %v", err)
		}
	}

	// Оставляем оригинальное название файла
	fileName := fmt.Sprintf("img/reaction_%d_%s", id, fileHeader.Filename)

	// Открываем файл
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()

	// Загружаем в MinIO
	_, err = r.minioClient.PutObject(
		context.Background(),
		r.bucketName,
		fileName,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return fmt.Errorf("ошибка загрузки в MinIO: %w", err)
	}

	// Обновляем путь к изображению в базе
	reaction.Src = "http://localhost:9000/aspirinimages/" + fileName
	err = r.db.Save(&reaction).Error
	if err != nil {
		// Если не удалось сохранить в БД, удаляем из MinIO
		r.minioClient.RemoveObject(context.Background(), r.bucketName, fileName, minio.RemoveObjectOptions{})
		return fmt.Errorf("ошибка сохранения пути к изображению: %w", err)
	}

	return nil
}
