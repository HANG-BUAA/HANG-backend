package dao

import (
	"HANG-backend/src/custom_error"
	"HANG-backend/src/model"
	"HANG-backend/src/service/dto"
	"HANG-backend/src/utils"
	"fmt"
	"gorm.io/gorm"
)

var courseDao *CourseDao

type CourseDao struct {
	BaseDao
}

func NewCourseDao() *CourseDao {
	if courseDao == nil {
		courseDao = &CourseDao{
			NewBaseDao(),
		}
	}
	return courseDao
}

func (m *CourseDao) ConvertCourseModelsToOverviewDTOs(courses []model.Course) ([]dto.CourseOverviewDTO, error) {
	res := make([]dto.CourseOverviewDTO, 0)
	for _, course := range courses {
		tmp, err := m.ConvertCourseModelToOverviewDTO(&course)
		if err != nil {
			return nil, err
		}
		res = append(res, *tmp)
	}
	return res, nil
}

func (m *CourseDao) ConvertReviewModelsToOverviewDTOs(reviews []model.CourseReview, user *model.User) ([]dto.CourseReviewOverviewDTO, error) {
	res := make([]dto.CourseReviewOverviewDTO, 0)
	for _, review := range reviews {
		tmp, err := m.ConvertReviewModelToOverviewDTO(&review, user)
		if err != nil {
			return nil, err
		}
		res = append(res, *tmp)
	}
	return res, nil
}

func (m *CourseDao) ConvertMaterialModelsToOverviews(materials []model.CourseMaterial, user *model.User) ([]dto.CourseMaterialOverviewDTO, error) {
	res := make([]dto.CourseMaterialOverviewDTO, 0)
	for _, material := range materials {
		tmp, err := m.ConvertMaterialModelToOverviewDTO(&material, user)
		if err != nil {
			return nil, err
		}
		res = append(res, *tmp)
	}
	return res, nil
}

func (m *CourseDao) ConvertCourseModelToOverviewDTO(course *model.Course) (*dto.CourseOverviewDTO, error) {
	// 查找 reviewNum
	var reviewNum int64
	if err := m.Orm.Model(&model.CourseReview{}).
		Where("course_id = ?", course.ID).
		Count(&reviewNum).
		Error; err != nil {
		return nil, err
	}

	// 查找平均分
	var average *float64
	if err := m.Orm.Model(&model.CourseReview{}).
		Where("course_id = ?", course.ID).
		Select("AVG(score)").
		Scan(&average).
		Error; err != nil {
		return nil, err
	}

	// 查询标签列表
	var tags []model.Tag
	if err := m.Orm.Table("course_tag").
		Select("tag.id, tag.name, tag.type").
		Joins("JOIN tag ON course_tag.tag_id = tag.id").
		Where("course_tag.course_id = ?", course.ID).
		Scan(&tags).
		Error; err != nil {
		return nil, err
	}
	return &dto.CourseOverviewDTO{
		ID:           course.ID,
		Name:         course.Name,
		Credits:      course.Credits,
		Campus:       course.Campus,
		ReviewNum:    int(reviewNum),
		AverageScore: average,
		Tags:         tags,
		CreatedAt:    course.CreatedAt,
		UpdatedAt:    course.UpdatedAt,
		DeletedAt:    course.DeletedAt,
	}, nil
}

func (m *CourseDao) ConvertReviewModelToOverviewDTO(review *model.CourseReview, user *model.User) (*dto.CourseReviewOverviewDTO, error) {
	isSelf := review.UserID == user.ID

	// 查找用户是否喜欢了
	hasLiked := true
	if err := m.Orm.Model(&model.CourseReviewLike{}).
		Where("course_review_id = ? AND user_id = ?", review.ID, user.ID).
		First(&model.CourseReviewLike{}).
		Error; err != nil {
		hasLiked = false
	}

	return &dto.CourseReviewOverviewDTO{
		ID:        review.ID,
		CourseID:  review.CourseID,
		Content:   review.Content,
		Score:     review.Score,
		IsSelf:    isSelf,
		HasLiked:  hasLiked,
		LikeNum:   review.LikeNum,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
		DeletedAt: review.DeletedAt,
	}, nil
}

func (m *CourseDao) ConvertMaterialModelToOverviewDTO(material *model.CourseMaterial, user *model.User) (*dto.CourseMaterialOverviewDTO, error) {
	isSelf := material.UserID == user.ID

	// 查找是否喜欢
	hasLiked := true
	if err := m.Orm.Model(&model.CourseMaterialLike{}).
		Where("course_material_id = ? AND user_id = ?", material.ID, user.ID).
		First(&model.CourseMaterialLike{}).
		Error; err != nil {
		hasLiked = false
	}

	res := dto.CourseMaterialOverviewDTO{
		ID:          material.ID,
		CourseID:    material.CourseID,
		Link:        material.Link,
		Source:      material.Source,
		Description: material.Description,
		IsOfficial:  material.IsOfficial,
		IsApproved:  material.IsApproved,
		LikeNum:     material.LikeNum,
		HasLiked:    hasLiked,
		IsSelf:      isSelf,
		CreatedAt:   material.CreatedAt,
		UpdatedAt:   material.UpdatedAt,
		DeletedAt:   material.DeletedAt,
	}

	// 如果非官方的，查找用户信息
	if !material.IsOfficial {
		var author model.User
		if err := m.Orm.Where("id = ?", material.UserID).First(&author).Error; err != nil {
			return nil, err
		}
		res.Author.ID = author.ID
		res.Author.Name = author.Username
		res.Author.Avatar = author.Avatar
	}
	return &res, nil
}

func (m *CourseDao) CreateCourse(id, name string, credits *float32, campus *int, tags []model.Tag) (*model.Course, error) {
	course := model.Course{
		ID:      id,
		Name:    name,
		Credits: credits,
		Campus:  campus,
	}

	// 使用数据库事务
	err := m.Orm.Transaction(func(tx *gorm.DB) error {
		// 1. 创建课程
		if err2 := tx.Create(&course).Error; err2 != nil {
			return err2
		}

		// 2. 加标签
		if tags != nil {
			for _, tag := range tags {
				courseTag := model.CourseTag{
					CourseID: course.ID,
					TagID:    tag.ID,
				}
				if err2 := tx.Create(&courseTag).Error; err2 != nil {
					return err2
				}
			}
		}
		return nil
	})
	return &course, err
}

func (m *CourseDao) CreateReview(courseID string, user *model.User, score int, content string) (*model.CourseReview, error) {
	courseReview := model.CourseReview{
		CourseID: courseID,
		UserID:   user.ID,
		Content:  content,
		Score:    score,
	}
	if err := m.Orm.Create(&courseReview).Error; err != nil {
		return nil, err
	}
	return &courseReview, nil
}

func (m *CourseDao) ListTagsByIDs(tags []uint) ([]model.Tag, error) {
	var result []model.Tag
	for _, tagID := range tags {
		var tag model.Tag
		if err := m.Orm.First(&tag, tagID).Error; err != nil {
			return nil, err
		}
		result = append(result, tag)
	}
	return result, nil
}

func (m *CourseDao) GetCourseByID(id string) (*model.Course, error) {
	var course model.Course
	if err := m.Orm.First(&course, model.Course{ID: id}).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (m *CourseDao) LikeReview(user *model.User, review *model.CourseReview) error {
	return m.Orm.Transaction(func(tx *gorm.DB) error {
		newReviewLike := model.CourseReviewLike{
			UserID:         user.ID,
			CourseReviewID: review.ID,
		}
		if err := tx.Create(&newReviewLike).Error; err != nil {
			return err
		}
		// 乐观锁动态维护
		result := tx.Model(&model.CourseReview{}).Where("id = ? AND like_version = ?", review.ID, review.LikeVersion).Updates(map[string]interface{}{
			"like_num":     review.LikeNum + 1,
			"like_version": review.LikeVersion + 1,
		})
		if result.RowsAffected == 0 {
			return custom_error.NewOptimisticLockError()
		}
		return nil
	})
}

func (m *CourseDao) CheckReviewLiked(user *model.User, courseReview *model.CourseReview) bool {
	var courseReviewLike model.CourseReviewLike
	if err := m.Orm.Where("user_id = ? AND course_review_id = ?", user.ID, courseReview.ID).First(&courseReviewLike).Error; err != nil {
		return false
	}
	return true
}

func (m *CourseDao) CheckMaterialLiked(user *model.User, courseMaterial *model.CourseMaterial) bool {
	var courseMaterialLike model.CourseMaterialLike
	if err := m.Orm.Where("user_id = ? AND course_material_id = ?", user.ID, courseMaterial.ID).First(&courseMaterialLike).Error; err != nil {
		return false
	}
	return true
}

func (m *CourseDao) ListCourse(cursor string, pageSize int, keyword string, tags []uint) ([]model.Course, int, bool, error) {
	query := m.Orm.Model(&model.Course{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if tags != nil && len(tags) > 0 {
		for i, tag := range tags {
			alias := fmt.Sprintf("t%d", i) // 为每个标签创建唯一别名
			query = query.Joins(fmt.Sprintf("JOIN course_tag %s ON %s.course_id = course.id AND %s.tag_id = ?", alias, alias, alias), tag)
		}
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, false, err
	}

	// 多查一条出来，判断是否有下一页
	var courses []model.Course
	query = query.
		Limit(pageSize + 1).
		Order("id desc")
	if cursor != "" {
		query = query.Where("id < ?", cursor)
	}
	if err := query.Find(&courses).Error; err != nil {
		return nil, 0, false, err
	}
	isEnd := len(courses) < pageSize+1
	return courses[:utils.IfThenElse(isEnd, len(courses), pageSize).(int)], int(total), isEnd, nil
}

func (m *CourseDao) CommonListReview(cursor *struct {
	LikeNum int
	ID      uint
}, pageSize int, courseID string) ([]model.CourseReview, int, bool, error) {
	query := m.Orm.Model(&model.CourseReview{}).
		Where("course_id = ?", courseID)

	// 先计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, false, err
	}

	// 多查一条记录出来，判断是否有下一项
	var coursesReviews []model.CourseReview
	query = query.
		Limit(pageSize + 1).
		Order("like_num desc").
		Order("id desc")
	if cursor != nil {
		query = query.Where("like_num < ?", cursor.LikeNum).Or("like_num = ? AND id < ?", cursor.LikeNum, cursor.ID)
	}
	if err := query.Find(&coursesReviews).Error; err != nil {
		return nil, 0, false, err
	}
	isEnd := len(coursesReviews) < pageSize+1
	return coursesReviews[:utils.IfThenElse(isEnd, len(coursesReviews), pageSize).(int)], int(total), isEnd, nil
}

func (m *CourseDao) GetReviewsByIDs(ids []uint) ([]model.CourseReview, error) {
	var reviews []model.CourseReview
	if err := m.Orm.Where("id IN (?)", ids).Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

func (m *CourseDao) CreateMaterial(user *model.User, courseID, link, description string, source int, isApproved, isOfficial bool) (*model.CourseMaterial, error) {
	courseMaterial := model.CourseMaterial{
		CourseID:    courseID,
		UserID:      user.ID,
		Link:        link,
		Description: description,
		Source:      source,
		IsApproved:  isApproved,
		IsOfficial:  isOfficial,
	}
	if err := m.Orm.Create(&courseMaterial).Error; err != nil {
		return nil, err
	}
	return &courseMaterial, nil
}

func (m *CourseDao) LikeMaterial(user *model.User, material *model.CourseMaterial) error {
	return m.Orm.Transaction(func(tx *gorm.DB) error {
		newMaterialLike := model.CourseMaterialLike{
			UserID:           user.ID,
			CourseMaterialID: material.ID,
		}
		if err := tx.Create(&newMaterialLike).Error; err != nil {
			return err
		}

		// 动态维护
		// 首先获取当前记录，确保Source字段正确 ———— beforeSave 钩子的问题
		var existingMaterial model.CourseMaterial
		if err := tx.Where("id = ?", material.ID).First(&existingMaterial).Error; err != nil {
			return err // 处理错误，例如记录不存在的情况
		}

		// 确保使用有效的Source值
		existingMaterial.LikeNum += 1
		existingMaterial.LikeVersion += 1

		// 执行更新
		result := tx.Save(&existingMaterial)
		if result.RowsAffected == 0 {
			return custom_error.NewOptimisticLockError()
		}

		return nil
	})
}

func (m *CourseDao) ListMaterial(cursor *struct {
	LikeNum int
	ID      uint
}, pageSize int, courseID string, isOfficial bool) ([]model.CourseMaterial, int, bool, error) {
	query := m.Orm.Model(&model.CourseMaterial{}).
		Where("course_id = ? AND is_approved = ? AND is_official = ?", courseID, true, isOfficial)

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, false, err
	}

	// 多查一条，判断是否有下一条记录
	var courseMaterials []model.CourseMaterial
	query = query.
		Limit(pageSize + 1).
		Order("like_num desc").
		Order("id desc")
	if cursor != nil {
		query = query.Where("like_num < ?", cursor.LikeNum).Or("like_num = ? AND id < ?", cursor.ID, cursor.ID)
	}
	if err := query.Find(&courseMaterials).Error; err != nil {
		return nil, 0, false, err
	}
	isEnd := len(courseMaterials) < pageSize+1
	return courseMaterials[:utils.IfThenElse(isEnd, len(courseMaterials), pageSize).(int)], int(total), isEnd, nil
}
