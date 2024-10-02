import React, { useState } from 'react';

function Dictionary() {
    const [word, setWord] = useState('');
    const [result, setResult] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [audioSources, setAudioSources] = useState([]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setResult(null);
        try {
            const response = await fetch(`http://localhost:8080/word?text=${word}`);
            if (!response.ok) {
                throw new Error('Word not found');
            }
            const data = await response.json();
            setResult(data);
            // Generate audio sources
            if (data.audioPronunciations && data.audioPronunciations.length > 0) {
                const sources = data.audioPronunciations.map(audio =>
                    `https://media.merriam-webster.com/audio/prons/en/us/mp3/${audio[0]}/${audio}.mp3`
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

    return (
        <div className="container mx-auto p-6 vintage-background rounded-lg shadow-lg"> {/* Use the new class */}
            <h1 className="text-4xl font-bold mb-6 text-center text-blue-600">GoDictionary</h1>
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
            {loading && <p className="text-center text-blue-500">Loading...</p>}
            {error && <p className="text-center text-red-500">{error}</p>}
            {result && (
                <div className="bg-white p-4 rounded-lg shadow-md">
                    <h2 className="text-3xl font-bold mb-4 text-blue-800">{result.text}</h2>
                    {result.pronunciations && result.pronunciations.length > 0 && (
                        <p className="mb-2"><strong>Pronunciation:</strong> {result.pronunciations.join(', ')}</p>
                    )}
                    {result.ipaPronunciation && (
                        <p className="mb-2"><strong>IPA Pronunciation:</strong> /{result.ipaPronunciation}/</p>
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
                    {result.definitions && result.definitions.length > 0 && result.definitions.map((def, index) => (
                        <div key={index} className="mb-4">
                            <h3 className="text-xl font-semibold">{def.partOfSpeech}</h3>
                            {def.senses && def.senses.length > 0 && (
                                <ol className="list-decimal text-center"> {/* Centering for ordered list */}
                                    {def.senses.map((sense, i) => (
                                        <li key={i}>{sense}</li>
                                    ))}
                                </ol>
                            )}
                        </div>
                    ))}
                    {result.idioms && result.idioms.length > 0 && (
                        <div className="mb-4">
                            <h3 className="text-xl font-semibold">Idioms</h3>
                            {result.idioms.map((idiom, index) => (
                                <div key={index} className="mb-2">
                                    <h4 className="font-semibold">{idiom.phrase}</h4>
                                    {idiom.senses && idiom.senses.length > 0 && (
                                        <ul className="list-disc text-center"> {/* Centering for unordered list */}
                                            {idiom.senses.map((sense, i) => (
                                                <li key={i}>{sense}</li>
                                            ))}
                                        </ul>
                                    )}
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            )}
        </div>
    );
}

export default Dictionary;