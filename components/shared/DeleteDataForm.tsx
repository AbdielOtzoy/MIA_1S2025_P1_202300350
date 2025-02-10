"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

interface DeleteDataFormProps {
  data: number[]
  onDelete: (index: number) => void
  onCancel: () => void
}

export default function DeleteDataForm({ data, onDelete, onCancel }: DeleteDataFormProps) {
  const [selectedIndex, setSelectedIndex] = useState<string>("")

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (selectedIndex !== "") {
      onDelete(Number(selectedIndex))
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <h2 className="text-2xl font-bold text-center text-gray-800 mb-4">Delete Data</h2>
      <Select value={selectedIndex} onValueChange={setSelectedIndex}>
        <SelectTrigger className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent">
          <SelectValue placeholder="Select an item to delete" />
        </SelectTrigger>
        <SelectContent>
          {data.map((item, index) => (
            <SelectItem key={index} value={index.toString()}>
              {`Index ${index}: ${item}`}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      <div className="flex justify-between">
        <Button
          type="submit"
          disabled={selectedIndex === ""}
          className="bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded-lg transition duration-300 ease-in-out transform hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Delete
        </Button>
        <Button
          type="button"
          onClick={onCancel}
          className="bg-gray-500 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded-lg transition duration-300 ease-in-out transform hover:scale-105"
        >
          Cancel
        </Button>
      </div>
    </form>
  )
}

