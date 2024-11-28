import React, { useState, useEffect } from 'react';
import SavedWordsList from './SavedWordsList';

// Import environment variables (e.g., using dotenv in development)
const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080';
const AUDIO_BASE_URL = process.env.REACT_APP_AUDIO_BASE_URL || 'https://media.example.com/audio/prons/en/us/mp3';

function Dictionary() {
    const [word, setWord] = useState('');
    const [result, setResult] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [audioSources, setAudioSources] = useState([]);
    const [savedWords, setSavedWords] = useState([]);
    const [showSavedWords, setShowSavedWords] = useState(false);

    useEffect(() => {
        fetchSavedWords();
    }, []);

    const fetchSavedWords = async () => {
        try {
            const response = await fetch(`${API_BASE_URL}/saved-words`);
            if (!response.ok) {
                throw new Error('Error fetching saved words');
            }
            const data = await response.json();
            setSavedWords(data);
        } catch (error) {
            console.error('Error:', error);
            setError(error.message);
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setResult(null);
        try {
            const response = await fetch(`${API_BASE_URL}/word?text=${word}`);
            if (!response.ok) {
                throw new Error('Word not found');
            }
            const data = await response.json();
            setResult(data);

            // Generate audio sources
            if (data.audioPronunciations && data.audioPronunciations.length > 0) {
                const sources = data.audioPronunciations.map(audio =>
                    `${AUDIO_BASE_URL}/${audio[0]}/${audio}.mp3`
                );
                setAudioSources(sources);
            }
        } catch (error) {
            console.error('Error:', error);
            setError(error.message);
        } finally {
            setLoading(false);
        }
    };

    const handleSaveWord = async () => {
        try {
            const meanings = result.definitions.flatMap(def => def.senses);
            const pronunciations = result.pronunciations;

            const response = await fetch(`${API_BASE_URL}/save`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    word,
                    meanings,
                    pronunciations,
                }),
            });
            if (!response.ok) {
                throw new Error('Error saving word');
            }

            // Update the savedWords state
            const newSavedWord = {
                word,
                meanings,
                pronunciations,
                savedDate: new Date().toLocaleString(),
            };
            setSavedWords(prevSavedWords => [...prevSavedWords, newSavedWord]);
            alert('Word saved successfully!');
        } catch (error) {
            console.error('Error:', error);
            alert(error.message);
        }
    };

    const toggleSavedWords = () => {
        setShowSavedWords(!showSavedWords);
    };

    return (
        <div className="container mx-auto p-6 vintage-background rounded-lg shadow-lg">
            <h1 className="text-4xl font-bold mb-6 text-center text-blue-600">GoDictionary</h1>
            <button onClick={toggleSavedWords} className="mb-4 bg-blue-600 text-white p-2 rounded hover:bg-blue-700 transition duration-200">
                {showSavedWords ? 'Back to Dictionary' : 'View Saved Words'}
            </button>
            {showSavedWords ? (
                <SavedWordsList savedWords={savedWords} />
            ) : (
                <form onSubmit={handleSubmit} className="mb-6 flex justify-center">
                    <input
                        type="text"
                        value={word}
                        onChange={(e) => setWord(e.target.value)}
                        placeholder="Enter a word"
                        className="border border-gray-300 p-3 rounded-l-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                    <button type="submit" className="bg-blue-600 text-white p-3 rounded-r-lg hover:bg-blue-700 transition duration-200">Look up</button>
                </form>
            )}
            {loading && <p className="text-center text-blue-500">Loading...</p>}
            {error && <p className="text-center text-red-500">{error}</p>}
            {result && !showSavedWords && (
                <div className="bg-white p-4 rounded-lg shadow-md">
                    <h2 className="text-3xl font-bold mb-4 text-blue-800">{result.text}</h2>
                    <button onClick={handleSaveWord} className="bg-green-500 text-white p-2 rounded hover:bg-green-600">Save Word</button>
                    {result.pronunciations && (
                        <p className="mb-2"><strong>Pronunciation:</strong> {result.pronunciations.join(', ')}</p>
                    )}
                    {audioSources.length > 0 && (
                        <div className="mb-2">
                            <strong>Audio Pronunciations:</strong>
                            {audioSources.map((src, index) => (
                                <audio key={index} controls>
                                    <source src={src} type="audio/mpeg" />
                                    Your browser does not support the audio element.
                                </audio>
                            ))}
                        </div>
                    )}
                </div>
            )}
        </div>
    );
}

export default Dictionary;
