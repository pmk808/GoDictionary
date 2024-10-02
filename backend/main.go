// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

type WordInfo struct {
	Text                string       `json:"text"`
	Pronunciations      []string     `json:"pronunciations"`
	IPAPronunciation    string       `json:"ipaPronunciation"`
	AudioPronunciations []string     `json:"audioPronunciations"` // New field for audio URLs
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

const apiKey = "258c6fb5-91a2-4294-a3db-2b52a33ea647"
const apiURL = "https://dictionaryapi.com/api/v3/references/collegiate/json/%s?key=%s"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/word", getWord).Methods("GET")

	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	http.Handle("/", corsMiddleware(r))
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func cleanText(text string) string {
	// Remove {sx|word||} patterns
	re := regexp.MustCompile(`\{sx\|[^}]+\|\|?\}`)
	text = re.ReplaceAllString(text, "")

	// Remove {bc} patterns
	text = strings.ReplaceAll(text, "{bc}", "")

	// Remove {dx_def} and {/dx_def} tags
	text = strings.ReplaceAll(text, "{dx_def}", "")
	text = strings.ReplaceAll(text, "{/dx_def}", "")

	// Remove {dxt|...} patterns
	re = regexp.MustCompile(`\{dxt\|[^}]+\}`)
	text = re.ReplaceAllString(text, "")

	// Remove leading numbers and dots
	re = regexp.MustCompile(`^\d+\.\s*`)
	text = re.ReplaceAllString(text, "")

	// Trim spaces
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
		// Extract pronunciations and IPA
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
						// Extract audio pronunciation
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
