import React from 'react';
import { Loader2 } from 'lucide-react';
import Logo from '@/assets/icons/ic-logo-badge.svg';
import { GLOBAL_CONFIG } from '@/global-config';

export const LoadingScreen: React.FC = () => {
  return (
    <div className="min-h-screen flex items-center justify-center bg-background">
      <div className="flex flex-col items-center space-y-4">
        <div className="flex items-center space-x-2">
          <img src={Logo} alt="Logo" className="w-8 h-8" />
          <span className="text-xl font-semibold">{GLOBAL_CONFIG.appName}</span>
        </div>
        <div className="flex items-center space-x-2 text-muted-foreground">
          <Loader2 className="w-4 h-4 animate-spin" />
          <span>正在初始化应用...</span>
        </div>
      </div>
    </div>
  );
};
