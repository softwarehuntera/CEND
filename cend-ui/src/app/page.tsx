'use client';

import { useState } from 'react';

import { FaCircleUser } from "react-icons/fa6";

import { SearchBar } from '@/components/search';
import { Filter } from '@/components/filter';
import { SearchResults } from '@/components/results';
import { DataSwap } from '@/components/swap';

export default function Home() {
  type SearchResult = {
    name: string;
    location: string;
    age: string;
    gender: string;
    [key: string]: string;
  };

  // Define four datasets
  const global: SearchResult[] = [
    { name: 'Roger', location: 'USA', age: '30', gender: 'Female' },
    { name: 'Roger', location: 'Canada', age: '25', gender: 'Male' },
    { name: 'Roger', location: 'USA', age: '35', gender: 'Male' },
    { name: 'Roger', location: 'USA', age: '28', gender: 'Female' },
  ];

  const dataset1: SearchResult[] = [
    { name: 'Alice', location: 'USA', age: '30', gender: 'Female' },
    { name: 'Bob', location: 'Canada', age: '25', gender: 'Male' },
    { name: 'Charlie', location: 'USA', age: '35', gender: 'Male' },
    { name: 'Diana', location: 'USA', age: '28', gender: 'Female' },
  ];

  const dataset2: SearchResult[] = [
    { name: 'Eve', location: 'Germany', age: '40', gender: 'Female' },
    { name: 'Frank', location: 'Italy', age: '50', gender: 'Male' },
    { name: 'Grace', location: 'Germany', age: '45', gender: 'Female' },
    { name: 'Hank', location: 'Italy', age: '55', gender: 'Male' },
  ];

  const dataset3: SearchResult[] = [
    { name: 'Ivy', location: 'France', age: '32', gender: 'Female' },
    { name: 'Jack', location: 'Spain', age: '29', gender: 'Male' },
    { name: 'Karen', location: 'France', age: '31', gender: 'Female' },
    { name: 'Leo', location: 'Spain', age: '33', gender: 'Male' },
  ];

  // When we have users working:
  // 1. Validate user
  // 2. Grab their preexisting datasets and update the datasets state

  // Store datasets in state    
  const [datasets] = useState([global, dataset1, dataset2, dataset3]);
  const [currentDatasetIndex, setCurrentDatasetIndex] = useState(0);

  // Get current dataset based on index
  const currentDataset = datasets[currentDatasetIndex];
  const [searchResults, setSearchResults] = useState(currentDataset);
  const [filteredResults, setFilteredResults] = useState(currentDataset);
  const [filters, setFilters] = useState<{ [key: string]: Set<string> }>({});

  const handleDatasetSwap = (index: number) => {
    setCurrentDatasetIndex(index);
    setSearchResults(datasets[index]);
    setFilteredResults(datasets[index]);
    setFilters({});
  };

  const handleFilterChange = (newFilters: { [key: string]: Set<string> }) => {
    setFilters(newFilters);
    applyFilters(newFilters);
  };

  const applyFilters = (activeFilters: { [key: string]: Set<string> }) => {
    if (Object.keys(activeFilters).length === 0) {
      setFilteredResults(searchResults);
      return;
    }

    const results = searchResults.filter((item) =>
      Object.entries(activeFilters).every(([field, values]) => {
        return values.has(item[field]);
      })
    );

    setFilteredResults(results);
  };

  return (
    <div className="grid min-h-screen p-4 gap-8 font-[family-name:var(--font-geist-sans)]">
     
      <main className="grid grid-cols-12 gap-4 w-full">

        <div className="col-span-10 flex flex-col gap-4">

          <SearchBar />

          <div className="grid grid-cols-12 gap-4 p-4 bg-white rounded-lg">

            <Filter searchResults={searchResults} onFilterChange={handleFilterChange} />
           
            <SearchResults results={filteredResults} />

          </div>
          
        </div>

        <div className="col-span-2 flex flex-col gap-4">

          <div className="flex justify-end">
            <FaCircleUser className="text-5xl text-white" />
          </div>

          <DataSwap onDatasetSwap={handleDatasetSwap} datasets={datasets} />

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
