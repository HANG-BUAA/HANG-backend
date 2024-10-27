package dao

import (
	"HANG-backend/src/model"
	"HANG-backend/src/service/dto"
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

func (m *CourseDao) ConvertCourseModelToOverviewDTO(course *model.Course) (*dto.CourseOverviewDTO, error) {
	// 查找 reviewNum
	var reviewNum int64
	if err := m.Orm.Model(&model.CourseReview{}).
		Where("course_id = ?", course.ID).
		Count(&reviewNum).
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
		ID:        course.ID,
		Name:      course.Name,
		Credits:   course.Credits,
		Campus:    course.Campus,
		ReviewNum: int(reviewNum),
		Tags:      tags,
		CreatedAt: course.CreatedAt,
		UpdatedAt: course.UpdatedAt,
		DeletedAt: course.DeletedAt,
	}, nil
}

func (m *CourseDao) ConvertReviewModelToOverviewDTO(review *model.CourseReview, user *model.User) (*dto.CourseReviewOverviewDTO, error) {
	isSelf := review.UserID == user.ID

	// 查找 likeNum
	var likeNum int64
	if err := m.Orm.Model(&model.CourseReviewLike{}).
		Where("course_review_id = ?", review.ID).
		Count(&likeNum).
		Error; err != nil {
		return nil, err
	}

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
		LikeNum:   int(likeNum),
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
		DeletedAt: review.DeletedAt,
	}, nil
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
	reviewLike := model.CourseReviewLike{
		UserID:         user.ID,
		CourseReviewID: review.ID,
	}
	if err := m.Orm.Create(&reviewLike).Error; err != nil {
		return err
	}
	return nil
}

func (m *CourseDao) CheckReviewLiked(user *model.User, courseReview *model.CourseReview) bool {
	var courseReviewLike model.CourseReviewLike
	if err := m.Orm.Where("user_id = ? AND course_review_id = ?", user.ID, courseReview.ID).First(&courseReviewLike).Error; err != nil {
		return false
	}
	return true
}
