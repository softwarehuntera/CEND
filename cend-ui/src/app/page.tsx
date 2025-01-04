'use client';

import { useEffect, useState } from 'react';

import { FaCircleUser } from "react-icons/fa6";

import { SearchBar } from '@/components/search';
import { Add } from '@/components/add';
import { TermSelection } from '@/components/select';
import { TermEntry } from '@/types/backendTypes';
import { SearchResults } from '@/components/results';
import { GetTerms } from '@/types/backendTypes';

export default function Home() {

  const [searchResults, setSearchResults] = useState<any[]>([]);
  const [selectedTerm, setSelectedTerm] = useState<TermEntry | null>(null);
  const [selectedTermCluster, setSelectedTermCluster] = useState<TermEntry[]>([]);

  const handleSearchResults = (results: any[]) => {
    setSearchResults(results);
  };

  const handleSelectedTerm = async (term: TermEntry) => {
    setSelectedTerm(term);
    const x = await getRelatedTerms(term);
    x? console.log(`Setting related terms ${x}`): console.log("No related terms found");
    x? setSelectedTermCluster(x): setSelectedTermCluster([]);
  };
  
    const getRelatedTerms = async (term: TermEntry) => {
      try {
        const getTerms: GetTerms = {ids: term.preferredDocuments};
        const response = await fetch("http://localhost:80/get", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          
          body: JSON.stringify(getTerms),
        });

        if (!response.ok) {
          throw new Error("Network response was not ok");
        }

        const data: TermEntry[] = await response.json();
        return data;
      } catch (error) {
        console.error("Error getting related terms:", error);
        alert("An error occurred while fetching the initial results.");
      }
    };

    // Function to fetch the initial search results on mount
    const onInitialization = async () => {
      try {
          const response = await fetch("http://localhost:80/query", {
              method: "POST",
              headers: {
                  "Content-Type": "application/json",
              },
              body: JSON.stringify({ min: 1, max: 100 }),
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
          {/* Results Section */}
          <SearchResults results={searchResults} onSelectTerm={handleSelectedTerm} />
        </div>

        {/* Right Column */}
        <div className="col-span-2 flex flex-col gap-4">
          {/* User Icon */}
          <div className="flex justify-end">
            <FaCircleUser className="text-5xl text-white" />
          </div>

          {/* Selected Term */}

          {/* Datasets Section */}
          {/* <div className="p-4 bg-gray-200 rounded-lg text-black">
            <h2 className="text-xl font-semibold">Datasets</h2>
            <p>Option 1</p>
            <p>Option 2</p>
            <p>Option 3</p>
          </div> */}

          {/* Add Section */}
          <TermSelection selectedTerm={selectedTerm} selectedTermCluster={selectedTermCluster} />
          <Add />
        </div>
      </main>
    </div>
  );
}
