import React from 'react';

interface StatusDisplayProps {
  status: string;
}

export const StatusDisplay: React.FC<StatusDisplayProps> = ({ status }) => {
  if (!status) return null;

  return (
    <div className="bg-white rounded-lg shadow-md p-4">
      <h3 className="text-sm font-medium text-gray-700 mb-2">Status:</h3>
      <div className="text-sm text-gray-600 font-mono bg-gray-50 p-3 rounded">
        {status}
      </div>
    </div>
  );
};