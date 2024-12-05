package api

import (
	"HANG-backend/src/global"
	"HANG-backend/src/model"
	"HANG-backend/src/service"
	"HANG-backend/src/service/dto"
	"HANG-backend/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"strconv"
)

type CourseApi struct {
	BaseApi
	Service *service.CourseService
}

func NewCourseApi() CourseApi {
	return CourseApi{
		BaseApi: NewBaseApi(),
		Service: service.NewCourseService(),
	}
}

func (m CourseApi) CreateCourse(c *gin.Context) {
	var requestDTO dto.AdminCourseCreateRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &requestDTO}).GetError(); err != nil {
		return
	}
	responseDTO, err := m.Service.CreateCourse(&requestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *responseDTO,
	})
}

func (m CourseApi) CreateReview(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	var requestDTO dto.CourseReviewCreateRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &requestDTO}).GetError(); err != nil {
		return
	}
	requestDTO.User = user
	responseDTO, err := m.Service.CreateReview(&requestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Data: *responseDTO,
	})
}

func (m CourseApi) LikeReview(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	review := c.MustGet("course_review").(*model.CourseReview)
	var requestDTO dto.CourseReviewLikeRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &requestDTO, BindParamsFromUri: true}).GetError(); err != nil {
		return
	}
	requestDTO.User = user
	requestDTO.CourseReview = review

	err := m.Service.LikeReview(&requestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: gin.H{
			"status": "like success",
		},
	})
}

func (m CourseApi) ListCourse(c *gin.Context) {
	cursor := c.MustGet("cursor").(string)
	pageSize := c.MustGet("page_size").(int)
	var requestDTO dto.CourseListRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &requestDTO}).GetError(); err != nil {
		return
	}
	requestDTO.Cursor = cursor
	requestDTO.PageSize = pageSize

	courses, err := m.Service.ListCourse(&requestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *courses,
	})
}

func (m CourseApi) UnlikeReview(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	review := c.MustGet("course_review").(*model.CourseReview)
	var requestDTO dto.CourseReviewUnlikeRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &requestDTO, BindParamsFromUri: true}).GetError(); err != nil {
		return
	}
	requestDTO.User = user
	requestDTO.CourseReview = review
	err := m.Service.UnlikeReview(&requestDTO)

	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: gin.H{
			"status": "unlike success",
		},
	})
}

func (m CourseApi) UnlikeMaterial(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	material := c.MustGet("course_material").(*model.CourseMaterial)

	var requestDTO dto.CourseMaterialUnlikeRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &requestDTO, BindParamsFromUri: true}).GetError(); err != nil {
		return
	}
	requestDTO.User = user
	requestDTO.CourseMaterial = material

	err := m.Service.UnlikeMaterial(&requestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: gin.H{
			"status": "unlike success",
		},
	})
}

func (m CourseApi) ListReview(c *gin.Context) {
	cursor := c.MustGet("cursor").(string)
	pageSize := c.MustGet("page_size").(int)
	user := c.MustGet("user").(*model.User)
	var requestDTO dto.CourseReviewListRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &requestDTO}).GetError(); err != nil {
		return
	}
	requestDTO.Cursor = cursor
	requestDTO.PageSize = pageSize
	requestDTO.User = user

	var reviews *dto.CourseReviewListResponseDTO
	var err error
	if requestDTO.Query != nil {
		reviews, err = m.Service.SearchListReview(&requestDTO)
	} else {
		reviews, err = m.Service.CommonListReview(&requestDTO)
	}
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *reviews,
	})
}

func (m CourseApi) Retrieve(c *gin.Context) {
	course := c.MustGet("course").(*model.Course)
	var requestDTO dto.CourseRetrieveRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &requestDTO}).GetError(); err != nil {
		return
	}
	requestDTO.Course = course

	courseOverview, err := m.Service.Retrieve(&requestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *courseOverview,
	})
}

func (m CourseApi) CreateMaterial(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	var requestDTO dto.CourseMaterialCreateRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &requestDTO}).GetError(); err != nil {
		return
	}
	requestDTO.User = user
	responseDTO, err := m.Service.CreateMaterial(&requestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *responseDTO,
	})
}

func (m CourseApi) LikeMaterial(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	material := c.MustGet("course_material").(*model.CourseMaterial)
	var requestDTO dto.CourseMaterialLikeRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &requestDTO, BindParamsFromUri: true}).GetError(); err != nil {
		return
	}
	requestDTO.User = user
	requestDTO.CourseMaterial = material

	err := m.Service.LikeMaterial(&requestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: gin.H{
			"status": "like success",
		},
	})
}

func (m CourseApi) ListMaterial(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	cursor := c.MustGet("cursor").(string)
	pageSize := c.MustGet("page_size").(int)
	var requestDTO dto.CourseMaterialListRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &requestDTO}).GetError(); err != nil {
		return
	}
	requestDTO.Cursor = cursor
	requestDTO.PageSize = pageSize
	requestDTO.User = user

	materials, err := m.Service.ListMaterial(&requestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Data: *materials,
	})

}

func (m CourseApi) ListTags(c *gin.Context) {
	type TagCountResponse struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	var tags []TagCountResponse
	m.Ctx = c

	// 使用gorm的LeftJoin查询，将标签与课程数量关联
	err := global.RDB.Table("tag"). // 确保表名是单数形式
					Select("tag.id as id, tag.name as name, COUNT(course_tag.course_id) as count").
					Joins("LEFT JOIN course_tag ON course_tag.tag_id = tag.id").
					Group("tag.id").
					Order("count DESC").
					Scan(&tags).Error
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Data: tags,
	})
}

func (m CourseApi) CreateCoursesByExcel(c *gin.Context) {
	type ExcelRow struct {
		ID      string   `json:"id"`
		Name    string   `json:"name"`
		Credits *float32 `json:"credits"`
		Campus  *int     `json:"campus"`
		Tags    []uint   `json:"tags"`
	}
	m.Ctx = c
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}

	f, err := excelize.OpenReader(file)
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}

	// 数据在第一个工作表中
	sheetName := f.GetSheetName(0)

	// 读取所有行数据（第一行是表头）
	rows, err := f.Rows(sheetName)
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}
	var data []ExcelRow

	// 跳过表头（假设第一行是表头）
	rowIndex := 0
	for rows.Next() {
		rowIndex++
		if rowIndex == 1 {
			// 第一行是表头，跳过
			continue
		}

		// 获取每一行的单元格数据
		columns, err := rows.Columns()
		if err != nil {
			m.Fail(ResponseJson{
				Msg: err.Error(),
			})
			return
		}

		// 解析数据
		credits := utils.ParseFloat32(columns[2])
		campus := utils.ParseInt(columns[3])
		tags := utils.ParseTags(columns[4])

		// 构造数据对象
		excelRow := ExcelRow{
			ID:      columns[0],
			Name:    columns[1],
			Credits: credits,
			Campus:  campus,
			Tags:    tags,
		}

		// 添加到结果集
		data = append(data, excelRow)
	}
	var res []dto.AdminCourseCreateResponseDTO
	for _, row := range data {
		requestDTO := dto.AdminCourseCreateRequestDTO{
			ID:      row.ID,
			Name:    row.Name,
			Credits: row.Credits,
			Campus:  row.Campus,
			Tags:    row.Tags,
		}
		responseDTO, err := m.Service.CreateCourse(&requestDTO)
		if err != nil {
			m.Fail(ResponseJson{
				Msg: err.Error(),
			})
			return
		}
		res = append(res, *responseDTO)
	}
	m.OK(ResponseJson{
		Data: res,
	})
}

func (m CourseApi) ApproveMaterial(c *gin.Context) {
	m.Ctx = c
	materialIDStr := c.Param("material_id") // 获取路径参数 material_id
	if materialIDStr == "" {
		m.Fail(ResponseJson{
			Msg: "material_id is required",
		})
		return
	}
	materialID, err := strconv.Atoi(materialIDStr)
	if err != nil || materialID <= 0 {
		m.Fail(ResponseJson{
			Msg: "material_id is invalid",
		})
		return
	}

	var material model.CourseMaterial
	if err := global.RDB.Where("id = ?", materialID).First(&material).Error; err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}

	// 判断 IsApproved 是否已经是 true
	if material.IsApproved {
		m.Fail(ResponseJson{
			Msg: "material is already approved",
		})
		return
	}

	// 更新 IsApproved 为 true
	material.IsApproved = true
	if err := global.RDB.Save(&material).Error; err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Data: material,
	})
}
