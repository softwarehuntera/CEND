'use client';

import { useState } from 'react';

import { FaCircleUser } from "react-icons/fa6";

import { SearchBar } from '@/components/search';
import { Filter } from '@/components/filter';
import { SearchResults } from '@/components/results';

export default function Home() {
  type SearchResult = {
    name: string;
    location: string;
    age: string;
    gender: string;
    [key: string]: string;
  };

  const mockResults: SearchResult[] = [
    { name: 'Alice', location: 'USA', age: '30', gender: 'Female' },
    { name: 'Bob', location: 'Canada', age: '25', gender: 'Male' },
    { name: 'Charlie', location: 'USA', age: '35', gender: 'Male' },
    { name: 'Diana', location: 'USA', age: '28', gender: 'Female' },
    // Add more mock results as needed
  ];

  const [searchResults, setSearchResults] = useState(mockResults);
  const [filteredResults, setFilteredResults] = useState(mockResults);
  const [filters, setFilters] = useState<{ [key: string]: string }>({});

  const handleFilterChange = (newFilters: { [key: string]: string }) => {
    setFilters(newFilters);
    applyFilters(newFilters);
  };

  const applyFilters = (activeFilters: { [key: string]: string }) => {
    if (Object.keys(activeFilters).length === 0) {
      setFilteredResults(searchResults);
      return;
    }

    const results = searchResults.filter((item) =>
      Object.entries(activeFilters).every(
        ([field, value]) => item[field] === value
      )
    );

    setFilteredResults(results);
  };

  return (
    <div className="grid min-h-screen p-4 gap-8 font-[family-name:var(--font-geist-sans)]">
      <main className="grid grid-cols-12 gap-4 w-full">
        {/* Left Column - Search and Results (10 cols) */}
        <div className="col-span-10 flex flex-col gap-4">
          {/* Search Input */}
          <SearchBar />

          {/* Results Container */}
          <div className="grid grid-cols-12 gap-4 p-4 bg-white rounded-lg">
            <Filter searchResults={searchResults} onFilterChange={handleFilterChange} />

            {/* Results Section */}
            <SearchResults results={filteredResults} />
          </div>
        </div>

        {/* Right Column - User Icon and Options (2 cols) */}
        <div className="col-span-2 flex flex-col gap-4">
          {/* User Icon */}
          <div className="flex justify-end">
            <FaCircleUser className="text-5xl text-white" />
          </div>

          {/* Datasets Section */}
          <div className="p-4 bg-gray-200 rounded-lg text-black">
            <h2 className="text-xl font-semibold">Datasets</h2>
            <p>Option 1</p>
            <p>Option 2</p>
            <p>Option 3</p>
          </div>

          {/* Add Section */}
          <div className="p-4 bg-gray-300 rounded-lg text-black">
            <h2 className="text-xl font-semibold">Add</h2>
            <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
              Add
            </button>
          </div>
        </div>
      </main>
    </div>
  );
}
