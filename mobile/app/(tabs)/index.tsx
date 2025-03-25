import React, { useState, useEffect, useCallback } from 'react';
import { ActivityIndicator, ScrollView, StyleSheet, View, FlatList } from 'react-native';
import { Stack } from 'expo-router';
import Header from '@/components/Header';
import ProductList from '@/components/ProductList';
import { CategoryType, ProductType } from '@/types/type';
import Categories from '@/components/Categories';
import FlashSale from '@/components/FlashSale';
import { getCategories, getSaleProducts, getProducts } from '@/service/ApiService';
import { Image } from "expo-image"
import productService, { SimpleProductObject } from '@/service/productService';
import categoryService, { Category } from '@/service/categoryService';


const HomeScreen = () => {
  const [products, setProducts] = useState<SimpleProductObject[]>([]);
  const [saleProducts, setSaleProducts] = useState<ProductType[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      const [categoryData, productData] = await Promise.all([
        categoryService.getAllCategories(),
        //getSaleProducts(),
        productService.getAllProducts(),
      ]);

      console.log("Dados recebidos: ", productData)
      setCategories(categoryData);
      //setSaleProducts(saleProductData);
      setProducts(productData);
    } catch (error) {
      if (__DEV__) {
        console.log(error);
      }
    } finally {
      setLoading(false);
    }
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
});

export default HomeScreen;