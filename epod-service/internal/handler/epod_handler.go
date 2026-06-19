package handler

import (
	"net/http"
	"strconv"

	"epod-service/internal/service"

	"github.com/gin-gonic/gin"
)

type EPODHandler struct {
	service *service.EPODService // Pastikan tipe ini sesuai dengan struct di service-mu
}

func NewEPODHandler(service *service.EPODService) *EPODHandler {
	return &EPODHandler{
		service: service,
	}
}

// 1. POST /upload (Menerima Form Data & File)
func (h *EPODHandler) Upload(c *gin.Context) {
	// Ambil data teks dari Form
	awb := c.PostForm("awb")
	receiver := c.PostForm("receiver")
	latStr := c.PostForm("lat")
	lngStr := c.PostForm("lng")

	// Validasi input teks dasar
	if awb == "" || receiver == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "AWB dan Nama Penerima wajib diisi"})
				return
	}

	// Konversi koordinat string ke float64
	lat, _ := strconv.ParseFloat(latStr, 64)
	lng, _ := strconv.ParseFloat(lngStr, 64)

	// Ambil file gambar "photo" dari request
	fileHeader, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Foto bukti pengiriman wajib diunggah"})
		return
	}

	// Buka file agar bisa dibaca sebagai io.Reader (bisa dioper ke service & storage)
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses file foto"})
		return
	}
	defer file.Close()

	// Oper data bersih ke Layer Service
	resp, err := h.service.ProcessUpload(c.Request.Context(), awb, receiver, lat, lng, file, fileHeader.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ePOD berhasil diunggah",
		"data":    resp,
	})
}

// 2. GET /epod (Mengambil semua data / filter status)
func (h *EPODHandler) List(c *gin.Context) {
	statusFilter := c.Query("status") // menangkap ?status=PENDING_REVIEW dari React

	epods, err := h.service.GetAll(c.Request.Context(), statusFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Format response disesuaikan dengan kebutuhan React: json.data || []
	c.JSON(http.StatusOK, gin.H{
		"data":  epods,
		"total": len(epods),
	})
}

// 3. GET /epod/:id
func (h *EPODHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	epod, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, epod)
}

// 4. GET /epod/awb/:awb
func (h *EPODHandler) GetByAWB(c *gin.Context) {
	awb := c.Param("awb")

	epod, err := h.service.GetByAWB(c.Request.Context(), awb)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, epod)
}

// 5. PATCH /epod/:id/verify (Menerima JSON status dari React)
func (h *EPODHandler) Verify(c *gin.Context) {
	id := c.Param("id")

	// Struct lokal untuk menangkap body JSON {"status": "VERIFIED"}
	var req struct {
		Status string `json:"status" binding:"required"`
	}

	// Karena React mengirim JSON untuk Verifikasi, di sini BARU boleh pakai BindJSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format status tidak valid"})
		return
	}

	err := h.service.VerifyEPOD(c.Request.Context(), id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status ePOD berhasil diperbarui"})
}