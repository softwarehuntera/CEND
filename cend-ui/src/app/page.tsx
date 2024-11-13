import Image from "next/image";
import { FaCircleUser } from "react-icons/fa6";
import { SearchBar } from '@/components/search';

export default function Home() {
  return (
    <div className="grid min-h-screen p-4 gap-8 font-[family-name:var(--font-geist-sans)]">
      <main className="grid grid-cols-12 gap-4 w-full">
        {/* Left Column - Search and Results (10 cols) */}
        <div className="col-span-10 flex flex-col gap-4">
          {/* Search Input */}
          <SearchBar/>

          {/* Results Container */}
          <div className="p-4 bg-white rounded-lg">
            <h1 className="text-xl font-semibold text-black">Results</h1>
            {/* Filter Section */}
            <div className="grid grid-cols-12 gap-4">
              <div className="col-span-1 p-4 bg-orange-500 rounded-lg text-black">
                <h2 className="text-xl font-semibold">Filter</h2>
                <p>Location</p>
                <p>Age</p>
                <p>Gender</p>
              </div>

              {/* Results Section */}
              <div className="col-span-11 p-4 bg-gray-500 rounded-lg text-black">
                <h2 className="text-xl font-semibold">Search Results</h2>
                {/* Map results here instead of static p tags */}
                {Array(8).fill(0).map((_, i) => (
                  <p key={i}>Result {i + 1}</p>
                ))}
              </div>
              </div>
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
            <p className="bg-orange-500">Option 1</p>
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
