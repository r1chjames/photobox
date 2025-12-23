package http

//
//import (
//	"github.com/gin-gonic/gin"
//	"gitlab.com/r1chjames/photobox/api/internal/core/service"
//	"net/http"
//)
//
//type PhotoHandler struct {
//	photoSvc service.PhotoService
//}
//
//func NewPhotoHandler(photoSvc service.PhotoService) PhotoHandler {
//	return PhotoHandler{photoSvc: photoSvc}
//}
//
//func (h *PhotoHandler) GetPhoto(c *gin.Context) {
//	id := c.Param("id")[1:]
//	photo, err := h.photoSvc.GetPhoto(id)
//	if err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"data": photo})
//}
//
//func (h *PhotoHandler) GetPhotoThumbnail(c *gin.Context) {
//	id := c.Param("id")[1:]
//	photo, err := h.photoSvc.GetPhotoThumbnail(id)
//	if err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
//		return
//	}
//	c.Data(http.StatusOK, "image/jpeg", photo)
//}
//
//func (h *PhotoHandler) GetPhotoBin(c *gin.Context) {
//	id := c.Param("id")[1:]
//	photo, err := h.photoSvc.GetPhotoBin(id)
//	if err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
//		return
//	}
//	c.Data(http.StatusOK, "image/jpeg", photo)
//}
//
//func (h *PhotoHandler) ListPhotos(c *gin.Context) {
//	photos, err := h.photoSvc.ListPhotos(c.Query("albumId"), c.Query("fromId"), c.Query("limit"), c.Query("thumbnail"))
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"data": photos})
//}
//
//func (h *PhotoHandler) GetPhotoCount(c *gin.Context) {
//	count, err := h.photoSvc.GetPhotoCount(c.Query("albumId"))
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"data": count})
//}
//
//func (h *PhotoHandler) IndexPhotos(c *gin.Context) {
//	err := h.photoSvc.IndexPhotos()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"message": "Indexing started"})
//}
