"use client";

import { useState } from "react";
import Terminal from "../components/Terminal";
import { Button } from "@/components/ui/button";
import { Upload } from "lucide-react";
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
  const handleExecute = async (command: string) => {
    if(command === "clear" || command === "cls") {
      setOutputContent([]);
      setInputContent([""]);
      return;
    }
    if (command.trim() !== "") {
      try {
        // Enviar el comando al servidor local usando la instancia de axios
        const response = await api.post("/execute", {
          command: command,
        });

        // Mostrar el comando ejecutado y la respuesta del servidor
        setOutputContent([...outputContent, `> ${command}`, response.data.message]);
      } catch (error) {
        // Manejar errores de la petición
        setOutputContent([...outputContent, `> ${command}`, "Error: Could not execute command"]);
        console.error("Error executing command:", error);
      }

      // Añadir una nueva línea vacía para seguir escribiendo
      setInputContent([...inputContent, ""]);
    }
  };

  const handleUpload = () => {
    setOutputContent([...outputContent, "File uploaded successfully"]);
  };

  return (
    <div className="min-h-screen p-8"
      style={{ fontFamily: "'Fira Code', monospace", backgroundColor: "#202020" }}
    >
      <h1 className="text-4xl font-bold mb-8 text-center text-white">File System EXT2</h1>

      <Terminal
        content={inputContent}
        contentOutput={outputContent}
        editable={true}
        onContentChange={handleContentChange}
        onExecute={handleExecute} 
        />

      <div className="flex justify-center space-x-4">
        <Button onClick={handleUpload} className="bg-blue-500 hover:bg-blue-600">
          <Upload className="mr-2 h-4 w-4" /> Upload File
        </Button>
      </div>
    </div>
  );
}