package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/muhadkprsnl/go-backend/internal/config"
	"github.com/muhadkprsnl/go-backend/internal/models"
	"github.com/muhadkprsnl/go-backend/internal/repositories"
	"github.com/muhadkprsnl/go-backend/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// Database and logger init
var client, _ = config.ConnectMongoDB()
var db *mongo.Database = client.Database("QA")
var logger = zap.NewExample()
var tableRepo = repositories.NewTableRepository(db, logger)

// --- Get all reports by environment ---
func GetReports(w http.ResponseWriter, r *http.Request) {
	env := r.URL.Query().Get("env")
	if env == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing environment query param")
		return
	}

	log.Println("🔍 Fetching reports for environment:", env)

	reports, err := tableRepo.GetReportsByEnvironment(env)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch reports")
		return
	}

	log.Printf("✅ Found %d reports for %s\n", len(reports), env)

	if reports == nil {
		reports = []models.FormData{}
	}
	utils.RespondWithJSON(w, http.StatusOK, reports)
}

// --- Update report by ID ---
// func UpdateReport(w http.ResponseWriter, r *http.Request) {
// 	id := r.URL.Query().Get("id")
// 	if id == "" {
// 		utils.RespondWithError(w, http.StatusBadRequest, "Missing ID")
// 		return
// 	}

// 	var updated models.FormData
// 	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
// 		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		utils.RespondWithError(w, http.StatusBadRequest, "Invalid ObjectID")
// 		return
// 	}

// 	updated.ObjectID = objID
// 	updated.CreatedAt = time.Now()

// 	if err := tableRepo.UpdateReport(objID, updated); err != nil {
// 		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update report")
// 		return
// 	}

// 	utils.RespondWithJSON(w, http.StatusOK, updated)
// }

// func UpdateReport(w http.ResponseWriter, r *http.Request) {
// 	id := r.URL.Query().Get("id")
// 	if id == "" {
// 		utils.RespondWithError(w, http.StatusBadRequest, "Missing ID")
// 		return
// 	}

// 	// // Log the incoming request body for debugging
// 	// bodyBytes, _ := io.ReadAll(r.Body)
// 	// r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset the body for decoding
// 	// logger.Info("Update request body", zap.String("body", string(bodyBytes)))

// 	// var updated models.FormData
// 	// if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
// 	// 	logger.Error("Failed to decode request body", zap.Error(err))
// 	// 	utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
// 	// 	return
// 	// }

// 	// Log the raw request body
// 	bodyBytes, _ := io.ReadAll(r.Body)
// 	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset for JSON decoding
// 	logger.Info("Raw request body", zap.String("body", string(bodyBytes)))

// 	var updated models.FormData
// 	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
// 		logger.Error("Failed to decode JSON", zap.Error(err))
// 		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		logger.Error("Invalid ObjectID", zap.String("id", id), zap.Error(err))
// 		utils.RespondWithError(w, http.StatusBadRequest, "Invalid ObjectID")
// 		return
// 	}

// 	updated.ObjectID = objID
// 	updated.CreatedAt = time.Now()

// 	if err := tableRepo.UpdateReport(objID, updated); err != nil {
// 		logger.Error("Failed to update report", zap.Error(err))
// 		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update report")
// 		return
// 	}

// 	utils.RespondWithJSON(w, http.StatusOK, updated)
// }

// --- Update report ---
// UpdateReport handles updating a report
func UpdateReport(w http.ResponseWriter, r *http.Request) {
	// Get ID from query parameter
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing id parameter")
		return
	}

	// Convert string ID to ObjectID
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// Parse request body
	var updateData models.FormData
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Update in repository
	if err := tableRepo.UpdateReport(id, updateData); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update report: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Report updated successfully"})
}

// --- Delete report ---
// func DeleteReport(w http.ResponseWriter, r *http.Request) {
// 	idParam := r.URL.Query().Get("id")
// 	if idParam == "" {
// 		utils.RespondWithError(w, http.StatusBadRequest, "Missing id parameter")
// 		return
// 	}

// 	id, err := primitive.ObjectIDFromHex(idParam)
// 	if err != nil {
// 		utils.RespondWithError(w, http.StatusBadRequest, "Invalid ID format")
// 		return
// 	}

// 	if err := tableRepo.DeleteReport(id); err != nil {
// 		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete report")
// 		return
// 	}

// 	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Report deleted successfully"})
// }

func DeleteReport(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL path
	vars := mux.Vars(r)
	idParam := vars["id"]
	if idParam == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing id parameter")
		return
	}

	// Convert string ID to ObjectID
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// Delete in repository
	if err := tableRepo.DeleteReport(id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete report")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Report deleted successfully"})
}
