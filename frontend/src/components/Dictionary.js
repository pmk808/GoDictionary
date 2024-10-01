import React, { useState } from 'react';

function Dictionary() {
    const [word, setWord] = useState('');
    const [definition, setDefinition] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await fetch(`http://localhost:8080/word?text=${word}`);
            const data = await response.json();
            setDefinition(data.definition);
        } catch (error) {
            console.error('Error:', error);
        }
    };

    return (
        <div>
            <h1>Dictionary</h1>
            <form onSubmit={handleSubmit}>
                <input
                    type="text"
                    value={word}
                    onChange={(e) => setWord(e.target.value)}
                    placeholder="Enter a word"
                />
                <button type="submit">Look up</button>
            </form>
            {definition && (
                <div>
                    <h2>{word}</h2>
                    <p>{definition}</p>
                </div>
            )}
        </div>
    );
}

export default Dictionary;