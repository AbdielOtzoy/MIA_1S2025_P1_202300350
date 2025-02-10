"use client"

import { Button } from "@/components/ui/button"

interface VisualDisplayProps {
  arrayData: number[]
  linkedListData: number[]
  doublyLinkedListData: number[]
  indexData: number[]
  structureType: string
  onBack: () => void
}

const colors = [
  "bg-red-400",
  "bg-blue-400",
  "bg-green-400",
  "bg-yellow-400",
  "bg-purple-400",
  "bg-pink-400",
  "bg-indigo-400",
  "bg-teal-400",
  "bg-cyan-400",
  "bg-rose-400",
]

export default function VisualDisplay({ arrayData, linkedListData, doublyLinkedListData, indexData, structureType, onBack }: VisualDisplayProps) {
  console.log(structureType)
  const data = structureType === "array" ? arrayData : structureType === "linkedlist" ? linkedListData : structureType === "DoublyLinkedList" ? doublyLinkedListData : indexData
  return (
    <div className="space-y-4">
      
      <div className="flex flex-wrap justify-center">
        {data.map((item, index) => (
          <div key={index} className="m-1">
            <div className={`w-10 h-10 flex items-center justify-center rounded-full ${colors[index % colors.length]}`}>
              {item}
            </div>
          </div>
        ))}
      </div>

      <div className="overflow-x-auto">
        <table className="min-w-full bg-white">
          <thead>
            <tr>
              <th className="py-2 px-4 border-b">Index</th>
              <th className="py-2 px-4 border-b">Value</th>
            </tr>
          </thead>
          <tbody>
            {data.map((item, index) => (
              <tr key={index}>
                <td className="py-2 px-4 border-b">{index}</td>
                <td className="py-2 px-4 border-b">{item}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      

      <Button
        type="button"
        onClick={onBack}
        className="bg-gray-500 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded-lg transition duration-300 ease-in-out transform hover:scale-105"
      >
        Back
      </Button>
    </div>
    
  )
}

