import { useState } from 'react';
import { UseStatusReturn } from '../types/window';

export const useStatus = (): UseStatusReturn => {
  const [status, setStatus] = useState<string>('');

  const clearStatus = () => setStatus('');

  return {
    status,
    setStatus,
    clearStatus,
  };
};