import React from 'react';

const SavedWordsList = ({ savedWords }) => {
    return (
        <div className="mt-6">
            <h2 className="text-2xl font-bold mb-4">Saved Words</h2>
            {savedWords.length > 0 ? (
                <ul className="list-disc pl-5">
                    {savedWords.map((savedWord, index) => (
                        <li key={index} className="mb-4 p-2 border-b border-gray-300">
                            <strong>{savedWord.word}</strong> - <em>{savedWord.savedDate}</em>
                            <br />
                            <span className="font-semibold">Meanings:</span> {savedWord.meanings.join(', ')}
                            <br />
                            <span className="font-semibold">Pronunciations:</span> {savedWord.pronunciations.join(', ')}
                        </li>
                    ))}
                </ul>
            ) : (
                <p>No saved words yet.</p>
            )}
        </div>
    );
};

export default SavedWordsList;