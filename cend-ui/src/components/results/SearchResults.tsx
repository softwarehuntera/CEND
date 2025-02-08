// src/components/results/SearchResults.tsx
import { TermEntry } from '@/types/backendTypes';
import React, { useState } from 'react';
import { Button } from '@/components/ui/button';

type SearchResultsProps = {
    onSelectTerm: (term: TermEntry) => void;
    results: TermEntry[];
    currentPage: number;
    handlePageChange: (pageNumber: number) => void;
};

export default function SearchResults({ results, onSelectTerm, currentPage, handlePageChange }: SearchResultsProps) {
    const resultsPerPage = 25;
    const totalPages = Math.ceil(results.length / resultsPerPage);
    const indexOfLastResult = currentPage * resultsPerPage;
    const indexOfFirstResult = indexOfLastResult - resultsPerPage;
    const currentResults = results.slice(indexOfFirstResult, indexOfLastResult);

    const getPageNumbers = () => Array.from({ length: totalPages }, (_, i) => i + 1);

    const excludedFields = ['secret'];
    console.log('Search results:', results);

    return (
        <div className="col-span-11 p-4 bg-gray-200 rounded-lg text-black">
            <h2 className="text-xl font-semibold">
                Search Results
                <span className="text-sm font-normal ml-2">
                    (Showing {indexOfFirstResult + 1}-{Math.min(indexOfLastResult, results.length)} of {results.length})
                </span>
            </h2>
                {results.length > 0 ? (
                    <>
                        <div className="space-y-2">
                            {currentResults.map((result, index) => (
                                <div
                                    key={index}
                                    className="p-2 bg-white rounded shadow cursor-pointer hover:bg-gray-100"
                                    onClick={() => onSelectTerm(result)}
                                >
                                    <p>
                                        <span className="font-semibold capitalize">Document: </span>
                                        {result.document}
                                    </p>
                                    <p>
                                        <span className="font-semibold capitalize">ID: </span>
                                        {result.id}
                                    </p>
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
                            ))}
                        </div>

                        {totalPages > 1 && (
                            <div className="flex justify-center items-center gap-2 mt-4">
                                <Button
                                    variant="outline"
                                    onClick={() => handlePageChange(currentPage - 1)}
                                    disabled={currentPage === 1}
                                    className="px-3 py-1"
                                >
                                    Previous
                                </Button>
                                
                                {getPageNumbers().map((number) => (
                                    <Button
                                        key={number}
                                        variant={currentPage === number ? "default" : "outline"}
                                        onClick={() => handlePageChange(number)}
                                        className="px-3 py-1"
                                    >
                                        {number}
                                    </Button>
                                ))}
                                
                                <Button
                                    variant="outline"
                                    onClick={() => handlePageChange(currentPage + 1)}
                                    disabled={currentPage === totalPages}
                                    className="px-3 py-1"
                                >
                                    Next
                                </Button>
                            </div>
                    )}
                </>
            ) : (
                <p>No results found.</p>
            )}
        </div>
    );
}