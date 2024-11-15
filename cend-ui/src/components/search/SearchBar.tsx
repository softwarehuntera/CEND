import { useState, ChangeEvent, KeyboardEvent } from "react";

export default function SearchBar() {
    const [inputValue, setInputValue] = useState<string>("");

    const handleInputChange = (e: ChangeEvent<HTMLInputElement>) => {
        setInputValue(e.target.value);
    };

    const handleSubmit = async () => {
        if (inputValue.trim() === "") {
            alert("Please enter a search term.");
            return;
        }

        try {
            const response = await fetch("http://localhost:8000/search", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ query: inputValue }), // matches SearchRequest struct in JSON format
            });

            if (!response.ok) {
                throw new Error("Network response was not ok");
            }

            const data = await response.json();
            console.log("Search results:", data); // handle the response as needed
        } catch (error) {
            console.error("Fetch error:", error);
            alert("An error occurred while fetching search results.");
        }
    };

    const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter") {
            handleSubmit();
        }
    };

    return (
        <div className="bg-white">
            <div className="flex space-x-2">
                <input
                    type="text"
                    placeholder="Search..."
                    value={inputValue}
                    onChange={handleInputChange}
                    onKeyDown={handleKeyDown}
                    className="w-full p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
                />
                <button
                    onClick={handleSubmit}
                    className="p-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                    Submit
                </button>
            </div>
        </div>
    );
}
