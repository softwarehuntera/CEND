import Image from "next/image";
import { FaCircleUser } from "react-icons/fa6";

export default function Home() {
  return (
    <div className="grid items-start min-h-screen p-4 gap-8 font-[family-name:var(--font-geist-sans)]">

      <main className="flex flex-col gap-8 items-start w-full">

        {/* Search Input and Icon Container */}
        <div className="flex w-full items-center gap-4">
          {/* Search Input */}
          <input
            type="text"
            placeholder="Search..."
            className="w-full p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
          />

          {/* User Icon on the Right Side */}
          <FaCircleUser className="text-5xl text-white ml-auto" />
        </div>

        <div className="grid grid-cols-12 gap-4 w-full">

          {/* Main Section */}
          <div className="col-span-10 p-4 bg-gray-100 rounded-lg text-black">
            <h2 className="text-xl font-semibold">Results</h2>

            <div className="grid grid-cols-12 gap-4 w-full">

              {/* Filter Section */}
              <div className="col-span-1 p-4 bg-orange-500 rounded-lg text-black">
                <h2 className="text-xl font-semibold">Filter</h2>
                <p>Location</p>
                <p>Age</p>
                <p>Gender</p>
              </div>

              {/* Data Results Section */}
              <div className="col-span-11 p-4 bg-gray-500 rounded-lg text-black">
                <h2 className="text-xl font-semibold">Search Results</h2>
                <p>Result</p>
                <p>Result</p>
                <p>Result</p>
                <p>Result</p>
                <p>Result</p>
                <p>Result</p>
                <p>Result</p>
                <p>Result</p>
              </div>

            </div>
          </div>

          {/* Sub Section with Two Stacked Divs */}
          <div className="col-span-2 flex flex-col gap-4">
            {/* First Stacked Div */}
            <div className="p-4 bg-gray-200 rounded-lg text-black">
              <h2 className="text-xl font-semibold">Datasets</h2>
              <p className="bg-orange-500">Option 1</p>
              <p>Option 2</p>
              <p>Option 3</p>
            </div>

            {/* Second Stacked Div */}
            <div className="p-4 bg-gray-300 rounded-lg text-black">
              <h2 className="text-xl font-semibold">Add</h2>
              <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
                Add
              </button>
            </div>
          </div>

        </div>

      </main>
    </div>
  );
}
