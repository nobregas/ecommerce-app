import { ActivityIndicator, Image, ScrollView, StyleSheet, View, } from 'react-native'
import { Stack } from 'expo-router'
import Header from '@/components/Header'

import ProductList from '@/components/ProductList'
import { CategoryType, ProductType } from '@/types/type'
import { useState, useEffect } from 'react'
import Categories from '@/components/Categories'
import FlashSale from '@/components/FlashSale'
import { getCategories, getSaleProducts, getProducts } from '@/service/ApiService'

type Props = {}

const HomeScreen = (props: Props) => {

  const [products, setProducts] = useState<ProductType[]>([])
  const [saleProducts, setSaleProducts] = useState<ProductType[]>([])
  const [categories, setCategories] = useState<CategoryType[]>([])
  const [loading, setLoading] = useState<boolean>(true)

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    try{
      const[categoryData, saleProductData, productData] = await Promise.all([
        getCategories(),
        getSaleProducts(),
        getProducts(),
      ])

      setCategories(categoryData)
      setSaleProducts(saleProductData)
      setProducts(productData)

    } catch (error) {
      console.log(error)
    } finally {
      setLoading(false)
    }
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