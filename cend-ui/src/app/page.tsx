'use client';

import { useEffect, useState } from 'react';

import { FaCircleUser } from "react-icons/fa6";

import { SearchBar } from '@/components/search';
import { Add } from '@/components/add';
// import { Filter } from '@/components/filter';
import { SearchResults } from '@/components/results';

export default function Home() {
  // type SearchResult = {
  //   name: string;
  //   location: string;
  //   age: string;
  //   gender: string;
  //   [key: string]: string;
  // };

  // const mockResults: SearchResult[] = [
  //   { name: 'Alice', location: 'USA', age: '30', gender: 'Female' },
  //   { name: 'Bob', location: 'Canada', age: '25', gender: 'Male' },
  //   { name: 'Charlie', location: 'USA', age: '35', gender: 'Male' },
  //   { name: 'Diana', location: 'USA', age: '28', gender: 'Female' },
  //   // Add more mock results as needed
  // ];

  // const [searchResults, setSearchResults] = useState(mockResults);
  // const [filteredResults, setFilteredResults] = useState(mockResults);
  // const [filters, setFilters] = useState<{ [key: string]: Set<string> }>({});

  // const handleFilterChange = (newFilters: { [key: string]: Set<string> }) => {
  //   setFilters(newFilters);
  //   applyFilters(newFilters);
  // };

  // const applyFilters = (activeFilters: { [key: string]: Set<string> }) => {
  //   if (Object.keys(activeFilters).length === 0) {
  //     setFilteredResults(searchResults);
  //     return;
  //   }

  //   const results = searchResults.filter((item) =>
  //     Object.entries(activeFilters).every(([field, values]) => {
  //       // Include item if the field is not in filters or if the item's value is in the selected values
  //       return values.has(item[field]);
  //     })
  //   );

  //   setFilteredResults(results);
  // };

  const [searchResults, setSearchResults] = useState<any[]>([]);

  const handleSearchResults = (results: any[]) => {
    setSearchResults(results);
  };
  
    // Function to fetch the initial search results on mount
    const onInitialization = async () => {
      try {
          const response = await fetch("http://localhost:80/query", {
              method: "POST",
              headers: {
                  "Content-Type": "application/json",
              },
              body: JSON.stringify({ min: 0, max: 100 }),
          });

          if (!response.ok) {
              throw new Error("Network response was not ok");
          }

          const data = await response.json();
          // Update the search results with the initial data
          setSearchResults(data);
      } catch (error) {
          console.error("Error fetching initial search results:", error);
          alert("An error occurred while fetching the initial results.");
      }
  };

  // Call onInitialization when the component mounts
  useEffect(() => {
      onInitialization();
  }, []);

  return (
    <div className="grid min-h-screen p-4 gap-8 font-[family-name:var(--font-geist-sans)]">
      <main className="grid grid-cols-5 gap-4 w-full">
        {/* Left Column - Search and Results (10 cols) */}
        <div className="col-span-3 flex flex-col gap-4">
          {/* Search Input */}
          <SearchBar onSearchResults={handleSearchResults}/>

          {/* Results Container */}
          {/* <div className="grid grid-cols-12 gap-4 p-4 bg-white rounded-lg"> */}
            {/* <Filter searchResults={searchResults} onFilterChange={handleFilterChange} /> */}

            {/* Results Section */}
            <SearchResults results={searchResults} />
          {/* </div> */}
        </div>

        {/* Right Column - User Icon and Options (2 cols) */}
        <div className="col-span-2 flex flex-col gap-4">
          {/* User Icon */}
          <div className="flex justify-end">
            <FaCircleUser className="text-5xl text-white" />
          </div>

          {/* Datasets Section */}
          {/* <div className="p-4 bg-gray-200 rounded-lg text-black">
            <h2 className="text-xl font-semibold">Datasets</h2>
            <p>Option 1</p>
            <p>Option 2</p>
            <p>Option 3</p>
          </div> */}

          {/* Add Section */}
          <Add />
        </div>
      </main>
    </div>
  );
}
