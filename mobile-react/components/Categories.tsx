import { View, Text, StyleSheet, TouchableOpacity, FlatList } from 'react-native'
import React from 'react'

import { Colors } from '@/constants/Colors'
import { Image } from "expo-image"
import { Category } from '@/service/categoryService'

type Props = {
    categories: Category[]
}

const Categories = ({ categories }: Props) => {
    return (
        <View style={styles.container}>
            <View style={styles.titleWrapper}>
                <Text style={styles.title}>Categories</Text>
                <TouchableOpacity>
                    <Text style={styles.titleBtn}>See All</Text>
                </TouchableOpacity>
            </View>
            <FlatList
                data={categories}
                horizontal
                showsHorizontalScrollIndicator={false}
                keyExtractor={(item) => item.id.toString()}
                contentContainerStyle={{ paddingHorizontal: 20 }}
                renderItem={({ item, index }) => (
                    <TouchableOpacity>
                        <View style={styles.item}>
                            <Image source={{ uri: item.imageUrl }} style={styles.itemImg} />
                            <Text
                                numberOfLines={2}
                                ellipsizeMode='tail'
                                style={{
                                    textAlign: "center",
                                    fontSize: 12,
                                    lineHeight: 16,
                                    paddingHorizontal: 2,
                                }}
                            >
                                {item.name}
                            </Text>
                        </View>
                    </TouchableOpacity>
                )}
            />
        </View>
    )
}

export default Categories

const styles = StyleSheet.create({
    container: {
        marginBottom: 20,
        marginTop: 10,
    },
    titleWrapper: {
        flexDirection: "row",
        justifyContent: "space-between",
        marginHorizontal: 20,
        marginBottom: 10,
    },
    title: {
        fontSize: 18,
        fontWeight: "600",
        letterSpacing: 0.6,
        color: Colors.black
    },
    titleBtn: {
        fontSize: 14,
        fontWeight: "500",
        letterSpacing: 0.6,
        color: Colors.black
    },
    item: {
        marginVertical: 10,
        gap: 5,
        alignItems: "center",
        marginRight: 15,
        minWidth: 100,
        maxWidth: 120,
    },
    itemImg: {
        width: 60,
        height: 60,
        borderRadius: 30,
        backgroundColor: Colors.lightGray
    }
});