import React from 'react';
import { ProcessInfo, WindowInfo } from '../../bindings/hptools';
import { WindowDimensions } from '../types/window';
import { INPUT_LIMITS, SIZE_PRESETS, DEFAULT_DIMENSIONS } from '../constants/window';
import { NumberInput } from './NumberInput';
import { SizePresets } from './SizePresets';

interface WindowControlsProps {
  selectedProcess: ProcessInfo;
  dimensions: WindowDimensions;
  currentWindowInfo: WindowInfo | null;
  loading: boolean;
  onDimensionsChange: (dimensions: Partial<WindowDimensions>) => void;
  onSetWindowSize: () => void;
  onSetWindowPosition: () => void;
  onGetWindowInfo: () => void;
}

export const WindowControls: React.FC<WindowControlsProps> = ({
  selectedProcess,
  dimensions,
  currentWindowInfo,
  loading,
  onDimensionsChange,
  onSetWindowSize,
  onSetWindowPosition,
  onGetWindowInfo,
}) => {
  const handleDimensionChange = (key: keyof WindowDimensions) => (value: number) => {
    onDimensionsChange({ [key]: value });
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6 mb-6">
      <h2 className="text-xl font-semibold mb-4">
        Control Window: {selectedProcess.imageName}
      </h2>

      {/* Current Window Info */}
      <div className="mb-6">
        <button
          onClick={onGetWindowInfo}
          disabled={loading}
          className="px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600 disabled:opacity-50 mb-3"
        >
          📏 Get Current Window Info
        </button>
        
        {currentWindowInfo && (
          <div className="bg-gray-100 p-3 rounded-md text-sm">
            <strong>Current Window:</strong> {currentWindowInfo.width}x{currentWindowInfo.height} 
            at position ({currentWindowInfo.x}, {currentWindowInfo.y})
          </div>
        )}
      </div>

      {/* Size Controls */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
        <NumberInput
          label="Width"
          value={dimensions.width}
          onChange={handleDimensionChange('width')}
          min={INPUT_LIMITS.width.min}
          max={INPUT_LIMITS.width.max}
          defaultValue={DEFAULT_DIMENSIONS.width}
        />
        
        <NumberInput
          label="Height"
          value={dimensions.height}
          onChange={handleDimensionChange('height')}
          min={INPUT_LIMITS.height.min}
          max={INPUT_LIMITS.height.max}
          defaultValue={DEFAULT_DIMENSIONS.height}
        />
        
        <NumberInput
          label="X Position"
          value={dimensions.x}
          onChange={handleDimensionChange('x')}
          min={INPUT_LIMITS.position.min}
          max={INPUT_LIMITS.position.max}
          defaultValue={DEFAULT_DIMENSIONS.x}
        />
        
        <NumberInput
          label="Y Position"
          value={dimensions.y}
          onChange={handleDimensionChange('y')}
          min={INPUT_LIMITS.position.min}
          max={INPUT_LIMITS.position.max}
          defaultValue={DEFAULT_DIMENSIONS.y}
        />
      </div>

      {/* Action Buttons */}
      <div className="flex gap-4">
        <button
          onClick={onSetWindowSize}
          disabled={loading}
          className="px-6 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 disabled:opacity-50"
        >
          📏 Resize Window
        </button>
        
        <button
          onClick={onSetWindowPosition}
          disabled={loading}
          className="px-6 py-2 bg-purple-500 text-white rounded-md hover:bg-purple-600 disabled:opacity-50"
        >
          🎯 Move & Resize
        </button>
      </div>

      {/* Quick Size Presets */}
      <SizePresets 
        presets={SIZE_PRESETS} 
        onPresetSelect={onDimensionsChange} 
      />
    </div>
  );
};