abstract class ListaBase {
    abstract insertar(elementos: number[]): void;
    abstract eliminar(): void;
    abstract mostrar(): void;
}

export class ListaContigua extends ListaBase {
    private lista: number[];
    private capacidad: number;
    
    constructor(capacidad: number = 10) {
        super();
        this.capacidad = capacidad;
        this.lista = [];
    }

    insertar(elementos: number[]): void {
        if (this.lista.length + elementos.length <= this.capacidad) {
            this.lista.push(...elementos);
        } else {
            console.log("Lista llena, no se pueden insertar más elementos.");
        }
    }

    eliminar(): void {
        this.lista = [];
    }

    mostrar(): void {
        console.log("Lista Contigua:", this.lista.join(" "));
    }
}

// ===================== LISTA LIGADA =====================
class NodoSimple {
    valor: number;
    siguiente: NodoSimple | null;

    constructor(valor: number) {
        this.valor = valor;
        this.siguiente = null;
    }
}

export class ListaLigada extends ListaBase {
    private cabeza: NodoSimple | null;
    
    constructor() {
        super();
        this.cabeza = null;
    }

    insertar(valores: number[]): void {
        for (let i = valores.length - 1; i >= 0; i--) {
            const nuevoNodo = new NodoSimple(valores[i]);
            nuevoNodo.siguiente = this.cabeza;
            this.cabeza = nuevoNodo;
        }
    }

    eliminar(): void {
        this.cabeza = null;
    }

    mostrar(): void {
        let actual = this.cabeza;
        let resultado = "Lista Ligada: ";
        while (actual) {
            resultado += actual.valor + " -> ";
            actual = actual.siguiente;
        }
        resultado += "NULL";
        console.log(resultado);
    }
}

// ===================== LISTA DOBLEMENTE LIGADA =====================
class NodoDoble {
    valor: number;
    siguiente: NodoDoble | null;
    anterior: NodoDoble | null;

    constructor(valor: number) {
        this.valor = valor;
        this.siguiente = null;
        this.anterior = null;
    }
}

export class ListaDoblementeLigada extends ListaBase {
    private cabeza: NodoDoble | null;
    
    constructor() {
        super();
        this.cabeza = null;
    }

    insertar(valores: number[]): void {
        for (let i = 0; i < valores.length; i++) {
            const nuevoNodo = new NodoDoble(valores[i]);
            nuevoNodo.siguiente = this.cabeza;
            if (this.cabeza) {
                this.cabeza.anterior = nuevoNodo;
            }
            this.cabeza = nuevoNodo;
        }
    }

    eliminar(): void {
        this.cabeza = null;
    }

    mostrar(): void {
        let actual = this.cabeza;
        let resultado = "Lista Doblemente Ligada: ";
        while (actual) {
            resultado += actual.valor + " <-> ";
            actual = actual.siguiente;
        }
        resultado += "NULL";
        console.log(resultado);
    }
}

// ===================== LISTA INDEXADA =====================
export class ListaIndexada extends ListaBase {
    private lista: number[];
    private indices: number[];
    private capacidad: number;
    
    constructor(capacidad: number = 10) {
        super();
        this.capacidad = capacidad;
        this.lista = [];
        this.indices = [];
    }

    insertar(elementos: number[]): void {
        if (this.lista.length + elementos.length <= this.capacidad) {
            elementos.forEach((elemento) => {
                this.lista.push(elemento);
                this.indices.push(this.indices.length);
            });
        } else {
            console.log("Lista llena, no se pueden insertar más elementos.");
        }
    }

    eliminar(): void {
        this.lista = [];
        this.indices = [];
    }

    mostrar(): void {
        console.log("Lista Indexada:");
        this.indices.forEach((indice, i) => {
            console.log(`Índice: ${indice} -> Valor: ${this.lista[i]}`);
        });
    }
}

export type structureType = "array" | "linkedlist" | "DoublyLinkedList" | "index";