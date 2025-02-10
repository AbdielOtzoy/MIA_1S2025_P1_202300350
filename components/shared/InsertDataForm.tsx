"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"

import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { structureType } from "@/types"

interface InsertDataFormProps {
  onInsert: (value: number[], string: structureType) => void
  onCancel: () => void
}


export default function InsertDataForm({ onInsert, onCancel }: InsertDataFormProps) {
  const [values, setValues] = useState<number[]>([])
  const [structureType, setStructureType] = useState("")
  const [numElement, setNumElement] = useState(0);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (structureType !== "") {
      onInsert(values, structureType as structureType)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <h2 className="text-2xl font-bold text-center text-gray-800 mb-4">Insert Data for</h2>
      <Select value={structureType} onValueChange={setStructureType}>
        <SelectTrigger className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent">
          <SelectValue placeholder="Select an structure type" />
        </SelectTrigger>
        <SelectContent>
          {["array", "linkedlist", "DoublyLinkedList", "index"].map((item, index) => (
            <SelectItem key={index} value={item}>
              {item}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      <Input
        type="number"
        value={numElement}
        onChange={(e) => setNumElement(Number(e.target.value))}
        placeholder="Enter the number of elements"
        required
        max={10}
        min={1}
        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
      />
      <p className="text-md text-slate-400">Max elements: 10</p>
      {/* show inputs according the numElement variable and for input values*/}
      {[...Array(numElement)].map((_, index) => (
        <Input
          key={index}
          type="number"
          value={values[index] || ""}
          onChange={(e) => {
            const newValues = [...values]
            newValues[index] = Number(e.target.value)
            setValues(newValues)
          }}
          placeholder={`Enter element ${index + 1}`}
          required
          className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
      ))}
      <div className="flex justify-between">
        <Button
          type="submit"
          className="bg-green-500 hover:bg-green-600 text-white font-bold py-2 px-4 rounded-lg transition duration-300 ease-in-out transform hover:scale-105"
        >
          Insert
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

