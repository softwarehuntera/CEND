// src/components/results/SearchResults.tsx

type QueryResult = {
    document: string;
    id: number;
    fields: Record<string, string> | null;  // Fields can be a map of string keys to string values
    isPreferred: boolean;
    preferredDocuments: number[];
};

type SearchResultsProps = {
    results: QueryResult[];
};

export default function SearchResults({ results }: SearchResultsProps) {
    // Exclude these fields from being displayed
    const excludedFields = ['id', 'secret'];
    console.log('Search results:', results);

    return (
        <div className="col-span-11 p-4 bg-gray-200 rounded-lg text-black">
            <h2 className="text-xl font-semibold">Search Results</h2>
            {results.length > 0 ? (
                results.map((result, index) => (
                    <div key={index} className="p-2 bg-white my-2 rounded shadow">
                        {/* Display the document and id */}
                        <p>
                            <span className="font-semibold capitalize">Document: </span>
                            {result.document}
                        </p>
                        <p>
                            <span className="font-semibold capitalize">ID: </span>
                            {result.id}
                        </p>
                        {/* Display isPreferred and preferredDocuments */}
                        <p>
                            <span className="font-semibold capitalize">Is Preferred: </span>
                            {result.isPreferred ? 'Yes' : 'No'}
                        </p>
                        <p>
                            <span className="font-semibold capitalize">Preferred Documents: </span>
                            {result.preferredDocuments.length > 0 ? result.preferredDocuments.join(', ') : 'None'}
                        </p>
                        {result.fields && Object.entries(result.fields)
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