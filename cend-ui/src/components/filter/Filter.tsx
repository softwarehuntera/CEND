// src/components/filters/Filter.tsx
import { useState } from 'react';

interface FilterProps {
  searchResults: any[];
  onFilterChange: (filters: { [key: string]: Set<string> }) => void;
}

export default function Filter({ searchResults, onFilterChange }: FilterProps) {
  const [selectedFilters, setSelectedFilters] = useState<{ [key: string]: Set<string> }>({});

  // Extract unique fields and values
  const filterOptions = searchResults.reduce((acc, result) => {
    Object.keys(result).forEach((key) => {
      if (key !== 'name') {
        acc[key] = acc[key] || new Set();
        acc[key].add(result[key]);
      }
    });
    return acc;
  }, {} as { [key: string]: Set<string> });

  const handleFilterChange = (field: string, value: string) => {
    const newFilters = { ...selectedFilters };
    if (!newFilters[field]) newFilters[field] = new Set();

    if (value === 'ALL') {
      newFilters[field].clear();
      delete newFilters[field];
    } else {
      if (newFilters[field].has(value)) {
        newFilters[field].delete(value);
        if (newFilters[field].size === 0) {
          delete newFilters[field];
        }
      } else {
        newFilters[field].add(value);
      }
    }

    setSelectedFilters(newFilters);
    onFilterChange(newFilters); // Now called after state update
  };

  return (
    <div className="col-span-1 p-4 bg-orange-500 rounded-lg text-black">
      <h2 className="text-xl font-semibold">Filter</h2>
      {Object.entries(filterOptions).map(([field, values]) => (
        <div key={field}>
          <h3 className="font-semibold capitalize">{field}</h3>

          {/* "ALL" option */}
          <label className="block">
            <input
              type="checkbox"
              className="mr-2"
              checked={!selectedFilters[field]}
              onChange={() => handleFilterChange(field, 'ALL')}
            />
            All
          </label>

          {[...(values as Set<string>)].map((value) => (
            <label key={value} className="block">
              <input
                type="checkbox"
                className="mr-2"
                checked={selectedFilters[field]?.has(value) || false}
                onChange={() => handleFilterChange(field, value)}
              />
              {value}
            </label>
          ))}
        </div>
      ))}
    </div>
  );
}