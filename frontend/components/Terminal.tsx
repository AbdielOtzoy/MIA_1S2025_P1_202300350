import type React from "react";

interface TerminalProps {
  title: string;
  content: string[];
  editable?: boolean;
  onContentChange?: (index: number, newValue: string) => void;
  onExecute?: (command: string) => void; // Nueva prop para manejar la ejecución
}

const Terminal: React.FC<TerminalProps> = ({
  title,
  content,
  editable = false,
  onContentChange,
  onExecute,
}) => {
  const handleChange = (index: number) => (event: React.ChangeEvent<HTMLInputElement>) => {
    if (onContentChange) {
      onContentChange(index, event.target.value);
    }
  };

  const handleKeyDown = (index: number) => (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "Enter" && onExecute) {
      onExecute(content[index]); // Ejecuta el comando al presionar Enter
    }
  };

  return (
    <div className="flex flex-col h-full">
      <div className="bg-gray-800 text-white p-2 rounded-t-lg flex justify-between items-center">
        <div className="flex space-x-2 p-3">
          <div className="w-3 h-3 bg-red-500 rounded-full"></div>
          <div className="w-3 h-3 bg-yellow-500 rounded-full"></div>
          <div className="w-3 h-3 bg-green-500 rounded-full"></div>
        </div>
        <p className="font-semibold px-3 py-2">{title}</p>
      </div>
      <div className="flex-grow bg-black text-white p-4 rounded-b-lg overflow-auto font-mono max-h-96">
        {content.map((line, index) => (
          <div key={index} className="flex">
            <span className="text-gray-500 mr-4 select-none">{index + 1}</span>
            {editable ? (
              <input
                type="text"
                value={line || ""}
                onChange={handleChange(index)}
                onKeyDown={handleKeyDown(index)} // Maneja la tecla Enter
                className="bg-transparent outline-none flex-grow"
                autoFocus={index === content.length - 1} // Enfoca la última línea
              />
            ) : (
              <span>{line}</span>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};

export default Terminal;