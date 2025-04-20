import React, { useState, useEffect, useCallback } from 'react';
import { ActivityIndicator, ScrollView, StyleSheet, View, FlatList, TouchableOpacity } from 'react-native';
import { Stack, useRouter } from 'expo-router';
import Header from '@/components/Header';
import ProductList from '@/components/ProductList';
import { CategoryType, ProductType } from '@/types/type';
import Categories from '@/components/Categories';
import FlashSale from '@/components/FlashSale';
import { Image } from "expo-image"
import productService, { SimpleProductObject } from '@/service/productService';
import categoryService, { Category } from '@/service/categoryService';
import { useCartStore } from '@/store/cardBadgeStore';
import { Ionicons } from '@expo/vector-icons';
import { Colors } from '@/constants/Colors';

const HomeScreen = () => {
  const [products, setProducts] = useState<SimpleProductObject[]>([]);
  const [saleProducts, setSaleProducts] = useState<ProductType[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const { fetchCartItems } = useCartStore();
  const router = useRouter();

  useEffect(() => {
    fetchData();
    fetchCartItems();
  }, []);

  const fetchData = async () => {
    try {
      const [categoryData, productData] = await Promise.all([
        categoryService.getAllCategories(),
        productService.getAllProducts(),
      ]);

      setCategories(categoryData);
      setProducts(productData);
    } catch (error) {
      if (__DEV__) {
        console.log(error);
      }
    } finally {
      setLoading(false);
    }
  };

  const handleSupportChat = () => {
    router.push("/support/chat");
  };

  const renderCategories = useCallback(() => <Categories categories={categories} />, [categories]);
  const renderFlashSale = useCallback(() => <FlashSale products={saleProducts} />, [saleProducts]);
  const renderProductList = useCallback(() => <ProductList flatlist={false} products={products} />, [products]);

  if (loading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" />
      </View>
    );
  }

  return (
    <>
      <Stack.Screen options={{ headerShown: true, header: () => <Header /> }} />
      <ScrollView>
        {renderCategories()}
        {renderFlashSale()}
        <View style={styles.bannerWrapper}>
          <Image source={require("@/assets/images/sale-banner.jpg")} style={styles.banner} />
        </View>
        {renderProductList()}
      </ScrollView>

      {/* Floating Chat Support Button */}
      <TouchableOpacity 
        style={styles.chatButton} 
        onPress={handleSupportChat}
        activeOpacity={0.8}
      >
        <Ionicons name="chatbubbles" size={24} color={Colors.white} />
      </TouchableOpacity>
    </>
  );
};

const styles = StyleSheet.create({
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  bannerWrapper: {
    marginHorizontal: 20,
    marginBottom: 10,
  },
  banner: {
    width: "100%",
    height: 150,
    borderRadius: 15,
  },
  chatButton: {
    position: 'absolute',
    bottom: 20,
    right: 20,
    width: 56,
    height: 56,
    borderRadius: 28,
    backgroundColor: Colors.primary,
    justifyContent: 'center',
    alignItems: 'center',
    elevation: 5,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.25,
    shadowRadius: 3.84,
    zIndex: 1000,
  },
});

export default HomeScreen;