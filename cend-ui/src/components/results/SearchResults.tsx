// src/components/results/SearchResults.tsx

interface SearchResultsProps {
    results: any[];
}

export default function SearchResults({ results }: SearchResultsProps) {
    // Exclude these fields from being displayed
    const excludedFields = ['id', 'secret'];

    return (
        <div className="col-span-11 p-4 bg-gray-200 rounded-lg text-black">
            <h2 className="text-xl font-semibold">Search Results</h2>
            {results.length > 0 ? (
                results.map((result, index) => (
                    <div key={index} className="p-2 bg-white my-2 rounded shadow">
                        {Object.entries(result)
                            .filter(([key]) => !excludedFields.includes(key))
                            .map(([key, value]) => (
                                <p key={key}>
                                    <span className="font-semibold capitalize">{key}: </span>
                                    {value as React.ReactNode}
                                </p>
                            ))}
                    </div>
                ))
            ) : (
                <p>No results found.</p>
            )}
        </div>
    );
}