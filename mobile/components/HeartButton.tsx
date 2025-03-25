import { Ionicons } from "@expo/vector-icons";
import { TouchableOpacity } from "react-native";
import { Colors } from "@/constants/Colors";
import React, { useState, useEffect } from "react";
import productService from "@/service/productService";
import { useFavoritesStore } from "@/store/favoriteStore";
import { StyleSheet } from "react-native";

type HeartButtonProps = {
  productId: number;
  initialIsFavorited: boolean;
  size?: number;
  style?: any
};

const HeartButton = ({ productId, initialIsFavorited, size = 20, style }: HeartButtonProps) => {
    const [isFavorited, setIsFavorited] = useState(initialIsFavorited);
    const { toggleFavorite, isFavorite } = useFavoritesStore();
    const [isProcessing, setIsProcessing] = useState(false);
  
    useEffect(() => {
      const globalState = isFavorite(productId);
      if (globalState !== isFavorited) {
        setIsFavorited(globalState);
      }
    }, [productId, isFavorite]);
  
    const handlePress = async () => {
      if (isProcessing) return;
      
      try {
        setIsProcessing(true);
        const newState = !isFavorited;
        
        setIsFavorited(newState);
        toggleFavorite(productId, newState);
  
        if (newState) {
          await productService.addFavorite(productId);
        } else {
          await productService.removeFavorite(productId);
        }
      } catch (error) {
        const currentGlobalState = isFavorite(productId);
        setIsFavorited(currentGlobalState);
        alert("Erro ao atualizar favoritos");
        console.error("Erro no HeartButton:", error);
      } finally {
        setIsProcessing(false);
      }
    };

  return (
    <TouchableOpacity onPress={handlePress} style={style}>
      <Ionicons
        name={isFavorited ? "heart" : "heart-outline"}
        size={size}
        color={isFavorited ? Colors.primary : Colors.black}
      />
    </TouchableOpacity>
  );
};



export default HeartButton;