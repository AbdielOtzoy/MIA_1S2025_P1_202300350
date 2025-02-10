"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import InsertDataForm from "./InsertDataForm"
import VisualDisplay from "./VisualDisplay"
import { ListaContigua, ListaDoblementeLigada, ListaIndexada, ListaLigada, structureType } from "@/types"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

type MenuOption = "main" | "insert" | "delete" | "display" | "show"

export default function DataStructureVisualizer() {
  const [arrayData, setArrayData] = useState<number[]>([])
  const [linkedListData, setLinkedListData] = useState<number[]>([])
  const [doublyLinkedListData, setDoublyLinkedListData] = useState<number[]>([])
  const [indexData, setIndexData] = useState<number[]>([])
  const [currentMenu, setCurrentMenu] = useState<MenuOption>("main")
  const [selectedStructure, setSelectedStructure] = useState("array")

  const handleInsert = (value: number[], structure: structureType) => {
    switch (structure) {
      case "array":
        setArrayData(value)
        const newlista = new ListaContigua();
        newlista.insertar(value);
        newlista.mostrar();
        break
      case "linkedlist":
        setLinkedListData(value)
        const newLinkedList = new ListaLigada();
        newLinkedList.insertar(value);
        newLinkedList.mostrar();
        break
      case "DoublyLinkedList":
        setDoublyLinkedListData(value)
        const newDoublyLinkedList = new ListaDoblementeLigada();
        newDoublyLinkedList.insertar(value);
        newDoublyLinkedList.mostrar();
        break
      case "index":
        setIndexData(value)
        const newIndex = new ListaIndexada();
        newIndex.insertar(value);
        newIndex.mostrar();
        break
      }
        
    
    setCurrentMenu("main")
  }

  const handleDelete = (string: structureType) => {
    switch (string) {
      case "array":
        setArrayData([])
        break
      case "linkedlist":
        setLinkedListData([])
        break
      case "DoublyLinkedList":
        setDoublyLinkedListData([])
        break
      case "index":
        setIndexData([])
        break
    }
    setCurrentMenu("main")
  }


  const renderMenu = () => {
    switch (currentMenu) {
      case "insert":
        return <InsertDataForm onInsert={handleInsert} onCancel={() => setCurrentMenu("main")} />
      case "delete":
        return (
          <form onSubmit={(e) => {
            e.preventDefault()
            if (selectedStructure !== "") {
              handleDelete(selectedStructure as structureType)
            }
          }} className="space-y-6">
            <h2 className="text-2xl font-bold text-center text-gray-800 mb-4">Delete Data</h2>
            <Select value={selectedStructure} onValueChange={setSelectedStructure}>
              <SelectTrigger className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent">
                <SelectValue placeholder="Select an item to delete" />
              </SelectTrigger>
              <SelectContent>
                {/* show types of structures */}
                {["array", "linkedlist", "DoublyLinkedList", "index"].map((item, index) => (
                  <SelectItem key={index} value={item}>
                    {item}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <div className="flex justify-between">
              <Button
                type="submit"
                disabled={selectedStructure === ""}
                className="bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded-lg transition duration-300 ease-in-out transform hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Delete
              </Button>
              <Button
                type="button"
                onClick={() => setCurrentMenu("main")}
                className="bg-gray-500 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded-lg transition duration-300 ease-in-out transform hover:scale-105"
              >
                Cancel
              </Button>
            </div>
          </form>
        )
      case "display":
        return (
          <form onSubmit={(e) => {
            e.preventDefault()
            if (selectedStructure !== "") {
              setCurrentMenu("show")
            }
          }} className="space-y-6">
            <h2 className="text-2xl font-bold text-center text-gray-800 mb-4">Show Data</h2>
            <Select value={selectedStructure} onValueChange={setSelectedStructure}>
              <SelectTrigger className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent">
                <SelectValue placeholder="Select an item to delete" />
              </SelectTrigger>
              <SelectContent>
                {/* show types of structures */}
                {["array", "linkedlist", "DoublyLinkedList", "index"].map((item, index) => (
                  <SelectItem key={index} value={item}>
                    {item}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <div className="flex justify-between">
              <Button
                type="submit"
                disabled={selectedStructure === ""}
                className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded-lg transition duration-300 ease-in-out transform hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Show
              </Button>
              <Button
                type="button"
                onClick={() => setCurrentMenu("main")}
                className="bg-gray-500 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded-lg transition duration-300 ease-in-out transform hover:scale-105"
              >
                Cancel
              </Button>
            </div>
          </form>
        )
      case "show":
        return <VisualDisplay
                  arrayData={arrayData}
                  linkedListData={linkedListData}
                  doublyLinkedListData={doublyLinkedListData}
                  indexData={indexData}
                  structureType={selectedStructure as structureType}
                  onBack={() => setCurrentMenu("main")}
                />
      default:
        return (
          <div className="space-y-6">
            <Button
              onClick={() => setCurrentMenu("insert")}
              className="w-full bg-green-500 hover:bg-green-600 text-white font-bold py-3 rounded-lg transition duration-300 ease-in-out transform hover:scale-105"
            >
              1. Insert data
            </Button>
            <Button
              onClick={() => setCurrentMenu("delete")}
              className="w-full bg-red-500 hover:bg-red-600 text-white font-bold py-3 rounded-lg transition duration-300 ease-in-out transform hover:scale-105"
            >
              2. Delete data
            </Button>
            <Button
              onClick={() => setCurrentMenu("display")}
              className="w-full bg-blue-500 hover:bg-blue-600 text-white font-bold py-3 rounded-lg transition duration-300 ease-in-out transform hover:scale-105"
            >
              3. Display data visually
            </Button>
            <Button
              onClick={() => console.log("Exiting program")}
              className="w-full bg-gray-500 hover:bg-gray-600 text-white font-bold py-3 rounded-lg transition duration-300 ease-in-out transform hover:scale-105"
            >
              4. Exit the program
            </Button>
          </div>
        )
    }
  }

  return <div className="max-w-md mx-auto bg-white p-8 rounded-lg shadow-lg">{renderMenu()}</div>
}