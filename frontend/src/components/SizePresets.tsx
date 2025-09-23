import React from 'react';
import { SizePreset, WindowDimensions } from '../types/window';

interface SizePresetsProps {
  presets: readonly SizePreset[];
  onPresetSelect: (dimensions: Partial<WindowDimensions>) => void;
}

export const SizePresets: React.FC<SizePresetsProps> = ({
  presets,
  onPresetSelect,
}) => {
  const handlePresetClick = (preset: SizePreset) => {
    const dimensions: Partial<WindowDimensions> = {
      width: preset.w,
      height: preset.h,
    };
    
    if (preset.x !== undefined) {
      dimensions.x = preset.x;
    }
    if (preset.y !== undefined) {
      dimensions.y = preset.y;
    }
    
    onPresetSelect(dimensions);
  };

  return (
    <div className="mt-6">
      <h3 className="text-sm font-medium text-gray-700 mb-2">Quick Presets:</h3>
      <div className="flex gap-2 flex-wrap">
        {presets.map((preset) => (
          <button
            key={preset.name}
            onClick={() => handlePresetClick(preset)}
            className="px-3 py-1 text-xs bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
          >
            {preset.name} ({preset.w}x{preset.h}
            {preset.x !== undefined ? ` @${preset.x},${preset.y}` : ''})
          </button>
        ))}
      </div>
    </div>
  );
};