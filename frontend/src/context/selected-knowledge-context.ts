import { createContext, useContext } from "react";

import type { Document } from "@/types/document";
import type { KnowledgeSelection } from "@/types/knowledge";
import type { Note } from "@/types/note";

export interface SelectedKnowledgeContextType {
    selected: KnowledgeSelection | null;
    selectNote: (note: Note) => void;
    selectDocument: (document: Document) => void;
    clearSelection: () => void;
}

export const SelectedKnowledgeContext = createContext<
    SelectedKnowledgeContextType | undefined
>(undefined);

export function useSelectedKnowledge() {

    const context = useContext(
        SelectedKnowledgeContext,
    );

    if (!context) {
        throw new Error(
            "useSelectedKnowledge must be used inside SelectedKnowledgeProvider"
        );
    }

    return context;
}
