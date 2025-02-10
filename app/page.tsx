import DataStructureVisualizer from "@/components/shared/DataStructureVisualizer"

export default function Home() {
  return (
    <main className="min-h-screen bg-gradient-to-br from-purple-400 via-pink-500 to-red-500 p-8">
      <div className="container mx-auto bg-white bg-opacity-90 rounded-lg shadow-xl p-8">
        <h1 className="text-4xl font-bold mb-6 text-center text-gray-800">Data Structure Visualizer</h1>
        <DataStructureVisualizer />
      </div>
    </main>
  )
}