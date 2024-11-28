// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type WordInfo struct {
	Text                string       `json:"text"`
	Pronunciations      []string     `json:"pronunciations"`
	IPAPronunciation    string       `json:"ipaPronunciation"`
	AudioPronunciations []string     `json:"audioPronunciations"`
	Definitions         []Definition `json:"definitions"`
	Idioms              []Idiom      `json:"idioms"`
}

type Definition struct {
	PartOfSpeech string   `json:"partOfSpeech"`
	Senses       []string `json:"senses"`
}

type Idiom struct {
	Phrase string   `json:"phrase"`
	Senses []string `json:"senses"`
}

type SavedWord struct {
	Word           string   `json:"word"`
	Meanings       []string `json:"meanings"`
	Pronunciations []string `json:"pronunciations"`
	SavedDate      string   `json:"savedDate"`
}

const apiKey = "YOUR_API_KEY_HERE" // Replace with your actual API key
const apiURL = "https://dictionaryapi.com/api/v3/references/collegiate/json/%s?key=%s"
const savedWordsFile = "saved_words.json"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/word", getWord).Methods("GET")
	r.HandleFunc("/save", saveWordHandler).Methods("GET")
	r.HandleFunc("/saved-words", getSavedWordsHandler).Methods("GET")

	// Use the CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, 
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
	})

	handler := c.Handler(r)

	// Serve static files from the current directory
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("."))))

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func cleanText(text string) string {
	re := regexp.MustCompile(`\{sx\|[^}]+\|\|?\}`)
	text = re.ReplaceAllString(text, "")

	text = strings.ReplaceAll(text, "{bc}", "")
	text = strings.ReplaceAll(text, "{dx_def}", "")
	text = strings.ReplaceAll(text, "{/dx_def}", "")

	re = regexp.MustCompile(`\{dxt\|[^}]+\}`)
	text = re.ReplaceAllString(text, "")

	re = regexp.MustCompile(`^\d+\.\s*`)
	text = re.ReplaceAllString(text, "")

	return strings.TrimSpace(text)
}

func getWord(w http.ResponseWriter, r *http.Request) {
	word := r.URL.Query().Get("text")
	if word == "" {
		http.Error(w, "Missing 'text' parameter", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf(apiURL, word, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Error fetching data from API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading API response", http.StatusInternalServerError)
		return
	}

	var data []map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Error parsing API response", http.StatusInternalServerError)
		return
	}

	if len(data) == 0 {
		http.Error(w, "Word not found", http.StatusNotFound)
		return
	}

	wordInfo := WordInfo{Text: word}

	for _, entry := range data {
		if hwi, ok := entry["hwi"].(map[string]interface{}); ok {
			if prs, ok := hwi["prs"].([]interface{}); ok {
				for _, pr := range prs {
					if prMap, ok := pr.(map[string]interface{}); ok {
						if mw, ok := prMap["mw"].(string); ok {
							wordInfo.Pronunciations = append(wordInfo.Pronunciations, mw)
						}
						if ipa, ok := prMap["ipa"].(string); ok {
							wordInfo.IPAPronunciation = ipa
						}
						if sound, ok := prMap["sound"].(map[string]interface{}); ok {
							if audioID, ok := sound["audio"].(string); ok {
								wordInfo.AudioPronunciations = append(wordInfo.AudioPronunciations, audioID)
							}
						}
					}
				}
			}
		}

		if fl, ok := entry["fl"].(string); ok {
			def := Definition{PartOfSpeech: fl}
			if defs, ok := entry["def"].([]interface{}); ok {
				for _, d := range defs {
					if sseq, ok := d.(map[string]interface{})["sseq"].([]interface{}); ok {
						for _, ss := range sseq {
							if ssArray, ok := ss.([]interface{}); ok {
								for _, sense := range ssArray {
									if senseArray, ok := sense.([]interface{}); ok {
										for _, s := range senseArray {
											if sMap, ok := s.(map[string]interface{}); ok {
												if dt, ok := sMap["dt"].([]interface{}); ok {
													for _, t := range dt {
														if tArray, ok := t.([]interface{}); ok {
															if len(tArray) > 1 {
																if text, ok := tArray[1].(string); ok {
																	cleanedText := cleanText(text)
																	if cleanedText != "" {
																		def.Senses = append(def.Senses, cleanedText)
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
			if len(def.Senses) > 0 {
				wordInfo.Definitions = append(wordInfo.Definitions, def)
			}
		}

		if dros, ok := entry["dros"].([]interface{}); ok {
			for _, dro := range dros {
				if droMap, ok := dro.(map[string]interface{}); ok {
					if drp, ok := droMap["drp"].(string); ok {
						idiom := Idiom{Phrase: cleanText(drp)}
						if def, ok := droMap["def"].([]interface{}); ok {
							for _, d := range def {
								if sseq, ok := d.(map[string]interface{})["sseq"].([]interface{}); ok {
									for _, ss := range sseq {
										if ssArray, ok := ss.([]interface{}); ok {
											for _, sense := range ssArray {
												if senseArray, ok := sense.([]interface{}); ok {
													for _, s := range senseArray {
														if sMap, ok := s.(map[string]interface{}); ok {
															if dt, ok := sMap["dt"].([]interface{}); ok {
																for _, t := range dt {
																	if tArray, ok := t.([]interface{}); ok {
																		if len(tArray) > 1 {
																			if text, ok := tArray[1].(string); ok {
																				cleanedText := cleanText(text)
																				if cleanedText != "" {
																					idiom.Senses = append(idiom.Senses, cleanedText)
																				}
																			}
																		}
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
						if len(idiom.Senses) > 0 {
							wordInfo.Idioms = append(wordInfo.Idioms, idiom)
						}
					}
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wordInfo)
}

func saveWord(word string, meanings []string, pronunciations []string) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	savedWord := SavedWord{
		Word:           word,
		Meanings:       meanings,
		Pronunciations: pronunciations,
		SavedDate:      currentTime,
	}

	var savedWords []SavedWord
	file, err := os.Open(savedWordsFile)
	if err == nil {
		defer file.Close()
		bytes, _ := ioutil.ReadAll(file)
		json.Unmarshal(bytes, &savedWords)
	}

	savedWords = append(savedWords, savedWord)

	file, err = os.Create(savedWordsFile)
	if err != nil {
		return err
	}
	defer file.Close()
	json.NewEncoder(file).Encode(savedWords)

	return nil
}

func saveWordHandler(w http.ResponseWriter, r *http.Request) {
	word := r.URL.Query().Get("word")
	meanings := r.URL.Query()["meanings"]
	pronunciations := r.URL.Query()["pronunciations"]

	if word == "" {
		http.Error(w, "Missing 'word' parameter", http.StatusBadRequest)
		return
	}

	err := saveWord(word, meanings, pronunciations)
	if err != nil {
		http.Error(w, "Error saving word", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Word '%s' saved successfully", word)
}

func getSavedWordsHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(savedWordsFile)
	if err != nil {
		http.Error(w, "Error reading saved words file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var savedWords []SavedWord
	bytes, _ := ioutil.ReadAll(file)
	json.Unmarshal(bytes, &savedWords)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(savedWords)
}
