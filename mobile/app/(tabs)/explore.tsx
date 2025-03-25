import { FlatList, StyleSheet, Text, View } from 'react-native'
import React, { useEffect, useState } from 'react'
import { Stack } from 'expo-router'
import { useHeaderHeight } from '@react-navigation/elements'
import { Colors } from '@/constants/Colors'
import Animated, { FadeInDown } from 'react-native-reanimated'
import { Image } from "expo-image"
import categoryService, { Category } from '@/service/categoryService'

type Props = {}

const ExploreScreen = (props: Props) => {

  const headerHeight = useHeaderHeight()
  const [categories, setCategories] = useState<Category[]>([])

  useEffect(() => {
    fetchCategories()
  }, [])

  const fetchCategories = async () => {
    try {
      const categoriesData = await categoryService.getAllCategories()

      setCategories(categoriesData)

    } catch (error) {
      console.log(error)
    }
  }

  return (
    <>
      <Stack.Screen
        options={{ headerShown: true, headerTransparent: true, headerTitleAlign: 'center', title: 'explore' }}
      />

      <View style={[styles.container, { marginTop: headerHeight }]}>
        <FlatList
          data={categories}
          keyExtractor={(item) => item.id.toString()}
          showsVerticalScrollIndicator={false}
          renderItem={({ item, index }) => (
            <Animated.View style={styles.itemWrapper} entering={FadeInDown.delay(300 + index * 100).duration(500)}>
              <Text style={styles.itemTitle}>{item.name}</Text>
              <Image source={{ uri: item.imageUrl }} style={{ width: 100, height: 100, borderRadius: 10 }} />
            </Animated.View>
          )}
        />
      </View>
    </>
  )
}

export default ExploreScreen

const styles = StyleSheet.create({
  container: {
    flex: 1,
    paddingHorizontal: 20,
  },
  itemWrapper: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignContent: 'center',
    backgroundColor: Colors.extraLightGray,
    padding: 10,
    borderRadius: 10,
    marginBottom: 20
  },
  itemTitle: {
    fontSize: 16,
    fontWeight: '500',
    color: Colors.black
  },
})