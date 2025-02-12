"use client";

import { useState } from "react";
import Terminal from "../components/Terminal";
import { Button } from "@/components/ui/button";
import { Upload, Play, Trash2 } from "lucide-react";
import api from "@/lib/api";

export default function FileSystemSimulator() {
  const [inputContent, setInputContent] = useState<string[]>([""]);
  const [outputContent, setOutputContent] = useState<string[]>([]);

  // Función para manejar cambios en el contenido
  const handleContentChange = (index: number, newValue: string) => {
    const newContent = [...inputContent];
    newContent[index] = newValue;
    setInputContent(newContent);
  };

  // Función para manejar la ejecución de comandos
  const handleExecute = async () => {
    const lastCommand = inputContent[inputContent.length - 1]; // Obtiene el último comando
    if (lastCommand.trim() !== "") {
      try {
        // Enviar el comando al servidor local usando la instancia de axios
        const response = await api.post("/execute", {
          command: lastCommand,
        });

        // Mostrar el comando ejecutado y la respuesta del servidor
        setOutputContent([...outputContent, `> ${lastCommand}`, response.data.message]);
      } catch (error) {
        // Manejar errores de la petición
        setOutputContent([...outputContent, `> ${lastCommand}`, "Error: Could not execute command"]);
        console.error("Error executing command:", error);
      }

      // Añadir una nueva línea vacía para seguir escribiendo
      setInputContent([...inputContent, ""]);
    }
  };

  const handleUpload = () => {
    setOutputContent([...outputContent, "File uploaded successfully"]);
  };

  const handleClear = () => {
    setInputContent([""]);
    setOutputContent([]);
  };

  return (
    <div className="min-h-screen bg-gray-100 p-8"
      style={{ fontFamily: "'Fira Code', monospace", backgroundColor: "#202020" }}
    >
      <h1 className="text-4xl font-bold mb-8 text-center text-white">File System EXT2</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8">
        <Terminal
          title="Input Terminal"
          content={inputContent}
          editable={true}
          onContentChange={handleContentChange}
          onExecute={handleExecute} 
        />
        <Terminal title="Output Terminal" content={outputContent} />
      </div>
      <div className="flex justify-center space-x-4">
        <Button onClick={handleUpload} className="bg-blue-500 hover:bg-blue-600">
          <Upload className="mr-2 h-4 w-4" /> Upload File
        </Button>
        <Button onClick={handleExecute} className="bg-green-500 hover:bg-green-600">
          <Play className="mr-2 h-4 w-4" /> Execute
        </Button>
        <Button onClick={handleClear} className="bg-red-500 hover:bg-red-600">
          <Trash2 className="mr-2 h-4 w-4" /> Clear Console
        </Button>
      </div>
    </div>
  );
}