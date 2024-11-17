import React from "react";

type DataSwapProps = {
  onDatasetSwap: (index: number) => void;
  datasets: Array<any>; // Replace `any` with the actual dataset type if defined
};

export default function DataSwap({ onDatasetSwap, datasets }: DataSwapProps) {
  return (
    <div className="p-4 bg-gray-200 rounded-lg text-black">
      <h2 className="text-xl font-semibold">Datasets</h2>
      {/* <button
          key={0}
          onClick={() => onDatasetSwap(0)}
          className="block w-full mb-2 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        >
          Global
        </button> */}
      {datasets.map((_, index) => (
        <button
          key={index}
          onClick={() => onDatasetSwap(index)}
          className="block w-full mb-2 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        >
          Dataset {index + 1}
        </button>
      ))}
    </div>
  );
}
