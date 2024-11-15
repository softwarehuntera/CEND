// src/components/filters/Filter.tsx
import { useState } from 'react';

interface FilterProps {
  searchResults: any[];
  onFilterChange: (filters: { [key: string]: string }) => void;
}

export default function Filter({ searchResults, onFilterChange }: FilterProps) {
  const [selectedFilters, setSelectedFilters] = useState<{ [key: string]: string }>({});

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
    if (value === '') {
      delete newFilters[field]; // Deselect the filter
    } else {
      newFilters[field] = value; // Set the new selection
    }

    setSelectedFilters(newFilters);
    onFilterChange(newFilters);
  };

  return (
    <div className="col-span-1 p-4 bg-orange-500 rounded-lg text-black">
      <h2 className="text-xl font-semibold">Filter</h2>
      {Object.entries(filterOptions).map(([field, values]) => (
        <div key={field}>
          <h3 className="font-semibold capitalize">{field}</h3>

          {/* Optional "None" option */}
          <label key={`${field}-none`} className="block">
            <input
              type="radio"
              name={field}
              className="mr-2"
              checked={!selectedFilters[field]}
              onChange={() => handleFilterChange(field, '')}
            />
            None
          </label>

          {[...(values as Set<string>)].map((value) => (
            <label key={value} className="block">
              <input
                type="radio"
                name={field}
                className="mr-2"
                checked={selectedFilters[field] === value}
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