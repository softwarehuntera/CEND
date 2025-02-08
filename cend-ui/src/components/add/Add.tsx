// src/components/add/Add.tsx
import { NewTermEntry } from '@/types/backendTypes';
import { useState, FormEvent } from 'react';

interface DynamicField {
  key: string;
  value: string;
}

export default function Add() {
  const [name, setName] = useState('');
  const [isPreferred, setIsPreferred] = useState(false);
  const [dynamicFields, setDynamicFields] = useState<DynamicField[]>([]);
  const [newFieldKey, setNewFieldKey] = useState('');
  const [newFieldValue, setNewFieldValue] = useState('');
  const [preferredDocuments, setPreferredDocuments] = useState<number[]>([]);
  const [newDocumentId, setNewDocumentId] = useState('');

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    const termData: NewTermEntry = {
      document: name,
      isPreferred,
      preferredDocuments: preferredDocuments,
      fields: Object.fromEntries(dynamicFields.map(field => [field.key, field.value]))
    };

    try {
      const response = await fetch('http://localhost:80/add', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(termData),
      });
      console.log(response);
      if (!response.ok) {
        throw new Error('Failed to add term');
      }

      // Reset form
      setName('');
      setIsPreferred(false);
      setDynamicFields([]);
      setPreferredDocuments([]);
      alert('Term added successfully!');
    } catch (error) {
      console.error('Error adding term:', error);
      alert('Failed to add term');
    }
  };

  const addDynamicField = () => {
    if (newFieldKey && newFieldValue) {
      setDynamicFields([
        ...dynamicFields,
        { key: newFieldKey, value: newFieldValue }
      ]);
      setNewFieldKey('');
      setNewFieldValue('');
    }
  };

  const removeDynamicField = (index: number) => {
    setDynamicFields(dynamicFields.filter((_, i) => i !== index));
  };

  const addPreferredDocument = () => {
    const documentId = parseInt(newDocumentId, 10);
    if (!isNaN(documentId) && !preferredDocuments.includes(documentId)) {
      setPreferredDocuments([...preferredDocuments, documentId]);
      setNewDocumentId('');
    }
  };

  const removePreferredDocument = (index: number) => {
    setPreferredDocuments(preferredDocuments.filter((_, i) => i !== index));
  };

  return (
    <div className="p-4 bg-gray-200 rounded-lg shadow text-black">
      <h2 className="text-xl font-semibold mb-4">Add New Term</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        {/* Basic Fields */}
        <div>
          <label className="block mb-2">
            Term Name:
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full p-2 border rounded"
              required
            />
          </label>
        </div>

        <div>
          <label className="block mb-2">
            <input
              type="checkbox"
              checked={isPreferred}
              onChange={(e) => setIsPreferred(e.target.checked)}
              className="mr-2"
            />
            Preferred Term
          </label>
        </div>

        {/* Dynamic Fields Display */}
        <div className="space-y-2">
          {dynamicFields.map((field, index) => (
            <div key={index} className="flex items-center space-x-2">
              <span className="font-semibold">{field.key}:</span>
              <span>{field.value}</span>
              <button
                type="button"
                onClick={() => removeDynamicField(index)}
                className="text-red-500 hover:text-red-700"
              >
                Remove
              </button>
            </div>
          ))}
        </div>

        {/* Add New Field Section */}
        <div className="flex space-x-2">
          <input
            type="text"
            value={newFieldKey}
            onChange={(e) => setNewFieldKey(e.target.value)}
            placeholder="Field Name"
            className="p-2 border rounded"
          />
          <input
            type="text"
            value={newFieldValue}
            onChange={(e) => setNewFieldValue(e.target.value)}
            placeholder="Field Value"
            className="p-2 border rounded"
          />
          <button
            type="button"
            onClick={addDynamicField}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            Add Field
          </button>
        </div>

        {/* Preferred Documents */}
        <div className="space-y-2">
          <h3 className="font-semibold">Preferred Documents</h3>
          {preferredDocuments.map((docId, index) => (
            <div key={index} className="flex items-center space-x-2">
              <span>{docId}</span>
              <button
                type="button"
                onClick={() => removePreferredDocument(index)}
                className="text-red-500 hover:text-red-700"
              >
                Remove
              </button>
            </div>
          ))}
          <div className="flex space-x-2">
            <input
              type="text"
              value={newDocumentId}
              onChange={(e) => setNewDocumentId(e.target.value)}
              placeholder="Add Document ID"
              className="p-2 border rounded"
            />
            <button
              type="button"
              onClick={addPreferredDocument}
              className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              Add
            </button>
          </div>
        </div>

        {/* Submit Button */}
        <button
          type="submit"
          className="w-full px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
        >
          Add Term
        </button>
      </form>
    </div>
  );
}