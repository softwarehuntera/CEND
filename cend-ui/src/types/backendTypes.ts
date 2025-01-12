export type Query = {
    min: number;
    max: number;
};

export type TermEntry = {
    document: string;
    id: number;
    fields: Record<string, string> | null;
    isPreferred: boolean;
    preferredDocuments: number[];
};

export type NewTermEntry = {
    document: string;
    fields: Record<string, string> | null;
    isPreferred: boolean;
    preferredDocuments: number[];
};

export type GetTerms = {
    ids: number[];
}


