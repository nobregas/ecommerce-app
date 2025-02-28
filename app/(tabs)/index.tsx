import { ActivityIndicator, Image, ScrollView, StyleSheet, View, } from 'react-native'
import { Stack } from 'expo-router'
import Header from '@/components/Header'

import ProductList from '@/components/ProductList'
import { CategoryType, ProductType } from '@/types/type'
import axios from 'axios'
import { useState, useEffect } from 'react'
import Categories from '@/components/Categories'
import FlashSale from '@/components/FlashSale'

type Props = {}

const HomeScreen = (props: Props) => {

  const [products, setProducts] = useState<ProductType[]>([])
  const [saleProducts, setSaleProducts] = useState<ProductType[]>([])
  const [categories, setCategories] = useState<CategoryType[]>([])
  const [loading, setLoading] = useState<boolean>(true)

  useEffect(() => {
    getCategories()
    getSaleProducts()
    getProducts()
  }, [])

  const getProducts = async () => {
    const URL = `http://10.0.2.2:8000/products`
    const response = await axios.get(URL)

    setProducts(response.data)
    setLoading(false)
  }

  const getSaleProducts = async () => {
    const URL = `http://10.0.2.2:8000/saleProducts`
    const response = await axios.get(URL)

    setSaleProducts(response.data)
    setLoading(false)
  }

  const getCategories = async () => {
    const URL = `http://10.0.2.2:8000/categories`
    const response = await axios.get(URL)

    setCategories(response.data)
    setLoading(false)
  }

  if (loading) {
    return (
      <View>
        <ActivityIndicator size={"large"} />
      </View>
    )
  }

  return (
    <>
      <Stack.Screen options={{
        headerShown: true,
        header: () => <Header />
      }} />
      <ScrollView>
        <Categories categories={categories} />
        <FlashSale products={saleProducts} />
        <View style={styles.bannerWrapper}>
          <Image source={require("@/assets/images/sale-banner.jpg")} style={styles.banner} />
        </View>
        <ProductList flatlist={false} products={products} />
      </ScrollView>
    </>
  )
}

export default HomeScreen

const styles = StyleSheet.create({
  bannerWrapper: {
    marginHorizontal: 20,
    marginBottom: 10
  },
  banner: {
    width: "100%",
    height: 150,
    borderRadius: 15
  }

})